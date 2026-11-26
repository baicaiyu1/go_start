// Package stream 提供 SSE 流式推送示例：定时推送与「类大模型」逐段输出。
// 后续接大模型时，只替换 mock 文本来源，事件格式和断线处理不用改。
package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
)

const (
	// defaultDemoCount 未传 count 时默认推送的事件条数
	defaultDemoCount = 10
	// maxDemoCount count 上限，防止客户端要一个超长连接
	maxDemoCount = 100
	// minInterval 推送间隔下限，避免演示接口被当成压测工具
	minInterval = 100 * time.Millisecond
	// maxInterval 推送间隔上限，太慢会让客户端误以为连接卡死
	maxInterval = 10 * time.Second
	// defaultInterval 未传 interval_ms 时的默认推送间隔
	defaultInterval = time.Second
	// heartbeatEvery 心跳间隔：防止网关/NAT 60s 空闲断开
	heartbeatEvery = 15 * time.Second
	// chatChunkEvery 模拟大模型吐字间隔
	chatChunkEvery = 50 * time.Millisecond
)

// Handler 流式接口的 HTTP 层。
type Handler struct{}

// NewHandler 创建流式处理器。
func NewHandler() *Handler {
	return &Handler{}
}

// Demo GET /api/v1/stream/demo
//
// @Summary SSE 定时推送演示
// @Description 按 interval_ms 推送 count 条事件，每 15s 发一次心跳；客户端断开或写失败立即退出
// @Tags 流式
// @Produce text/event-stream
// @Param count query int false "事件条数，1-100" default(10)
// @Param interval_ms query int false "推送间隔（毫秒），100-10000" default(1000)
// @Success 200 {string} string "SSE 事件流"
// @Router /api/v1/stream/demo [get]
func (h *Handler) Demo(ctx context.Context, c *app.RequestContext) {
	count := clampInt(c.Query("count"), defaultDemoCount, 1, maxDemoCount)
	interval := clampMillis(c.Query("interval_ms"), defaultInterval, minInterval, maxInterval)

	w := sse.NewWriter(c)
	defer w.Close()

	// 首包立即推一条，客户端不用等到第一个 tick
	if err := writeEvent(w, "0", "open", utils.H{"message": "connected"}); err != nil {
		return
	}

	dataTicker := time.NewTicker(interval)
	defer dataTicker.Stop()
	keepAliveTicker := time.NewTicker(heartbeatEvery)
	defer keepAliveTicker.Stop()

	seq := 0
	for {
		select {
		case <-ctx.Done():
			// netpoll 且开启 WithSenseClientDisconnection 时客户端断开会让 ctx 取消；
			// 当前项目用 standard transporter，主要靠下面写失败来感知断开
			return
		case <-keepAliveTicker.C:
			if err := w.WriteKeepAlive(); err != nil {
				return
			}
		case <-dataTicker.C:
			seq++
			payload := utils.H{"seq": seq, "time": time.Now().Format(time.RFC3339)}
			if err := writeEvent(w, strconv.Itoa(seq), "message", payload); err != nil {
				return
			}
			if seq >= count {
				// done 事件让客户端知道可以停止重连
				_ = writeEvent(w, strconv.Itoa(seq), "done", utils.H{"message": "finished"})
				return
			}
		}
	}
}

// Chat POST /api/v1/stream/chat
//
// @Summary 流式对话（mock）
// @Description 接收 prompt，以 delta 事件逐段返回文本，最后发 done；接大模型时替换 mockReply 即可
// @Tags 流式
// @Accept json
// @Produce text/event-stream
// @Param req body ChatReq true "对话请求"
// @Success 200 {string} string "SSE 事件流"
// @Failure 400 {object} ErrorResp
// @Failure 422 {object} ErrorResp
// @Router /api/v1/stream/chat [post]
func (h *Handler) Chat(ctx context.Context, c *app.RequestContext) {
	var req ChatReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BindError(c, err)
		return
	}

	w := sse.NewWriter(c)
	defer w.Close()

	// 真实接入大模型时：把 mockReply 换成模型返回的 token 流（chan 或回调），
	// 每读到一段就 writeEvent(w, "", "delta", {"delta": "..."})，结束时发 done 即可
	for _, chunk := range chunkRunes(mockReply(req.Prompt), 4) {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := writeEvent(w, "", "delta", utils.H{"delta": chunk}); err != nil {
			return
		}
		time.Sleep(chatChunkEvery)
	}

	_ = writeEvent(w, "", "done", utils.H{"finish_reason": "stop"})
}

// writeEvent 把 data 序列化成 JSON 后写入一个 SSE 事件。
func writeEvent(w *sse.Writer, id, eventType string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteEvent(id, eventType, body)
}

// mockReply 伪造一段回复，仅用于验证流式链路。
func mockReply(prompt string) string {
	return fmt.Sprintf("收到你的问题：%s。当前是 mock 流式响应，接入大模型后替换这里即可。", prompt)
}

// chunkRunes 按固定长度切分文本，按 rune 切避免中文被截断成乱码。
func chunkRunes(s string, size int) []string {
	runes := []rune(s)
	chunks := make([]string, 0, (len(runes)+size-1)/size)
	for start := 0; start < len(runes); start += size {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

// clampInt 解析整数查询参数，非法或越界时收敛到默认值/边界。
func clampInt(raw string, def, minV, maxV int) int {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

// clampMillis 解析毫秒查询参数并转换为 time.Duration。
func clampMillis(raw string, def, minV, maxV time.Duration) time.Duration {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	d := time.Duration(v) * time.Millisecond
	if d < minV {
		return minV
	}
	if d > maxV {
		return maxV
	}
	return d
}
