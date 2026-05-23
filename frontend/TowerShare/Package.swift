// swift-tools-version: 5.10
import PackageDescription

let package = Package(
    name: "TowerShare",
    platforms: [
        .iOS(.v17),
        .macOS(.v14),
    ],
    products: [
        .library(name: "TowerShare", targets: ["TowerShare"]),
    ],
    dependencies: [
        // PocketBase Swift client
        .package(url: "https://github.com/pocketbase/swift-sdk", from: "0.2.0"),
    ],
    targets: [
        .target(
            name: "TowerShare",
            dependencies: [
                .product(name: "PocketBase", package: "swift-sdk"),
            ],
            path: "Sources/TowerShare"
        ),
        .testTarget(
            name: "TowerShareTests",
            dependencies: ["TowerShare"],
            path: "TowerShareTests"
        ),
    ]
)
