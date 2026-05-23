import PocketBase

// Single shared client. In production, read the URL from Info.plist / config.
let pb = PocketBase(url: URL(string: ProcessInfo.processInfo.environment["PB_URL"] ?? "http://localhost:8090")!)
