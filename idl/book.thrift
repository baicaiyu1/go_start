// 图书服务的 HTTP 接口定义（Hertz 使用的 Thrift IDL）
// 注意：hz 只识别 .thrift 后缀，使用前请把本文件改名为 book.thrift

namespace go book

// 图书实体：作为响应体返回
struct Book {
    1: i64 id (api.body = "id")
    2: string title (api.body = "title")
    3: string author (api.body = "author")
    4: double price (api.body = "price")
}

// GET /api/v1/books/:id 的请求：从 URL 路径取 id
struct GetBookReq {
    1: i64 id (api.path = "id")
}

// GET /api/v1/books 的请求：从 query 取分页参数
struct ListBooksReq {
    1: i32 page (api.query = "page")
    2: i32 page_size (api.query = "page_size")
}

// POST /api/v1/books 的请求：从 JSON 请求体取字段
struct CreateBookReq {
    1: string title (api.body = "title")
    2: string author (api.body = "author")
    3: double price (api.body = "price")
}

// 单本图书响应
struct BookResp {
    1: Book book (api.body = "book")
}

// 图书列表响应
struct BookListResp {
    1: list<Book> items (api.body = "items")
    2: i64 total (api.body = "total")
}

// 服务名会变成生成的 handler 文件名和路由分组依据
service BookService {
    BookResp GetBook(1: GetBookReq req) (api.get = "/api/v1/books/:id")
    BookListResp ListBooks(1: ListBooksReq req) (api.get = "/api/v1/books")
    BookResp CreateBook(1: CreateBookReq req) (api.post = "/api/v1/books")
}
