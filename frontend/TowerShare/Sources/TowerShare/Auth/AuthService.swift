import Foundation
import PocketBase

@MainActor
final class AuthService: ObservableObject {
    @Published var currentUser: RecordModel? = nil
    @Published var isLoading = false
    @Published var error: String? = nil

    static let shared = AuthService()

    private init() {
        currentUser = pb.authStore.model
    }

    func signUp(email: String, password: String, displayName: String) async {
        isLoading = true
        error = nil
        defer { isLoading = false }
        do {
            _ = try await pb.collection("users").create(body: [
                "email": email,
                "password": password,
                "passwordConfirm": password,
                "display_name": displayName,
            ])
            await signIn(email: email, password: password)
        } catch {
            self.error = error.localizedDescription
        }
    }

    func signIn(email: String, password: String) async {
        isLoading = true
        error = nil
        defer { isLoading = false }
        do {
            let result = try await pb.collection("users").authWithPassword(email, password)
            currentUser = result.record
        } catch {
            self.error = error.localizedDescription
        }
    }

    func signOut() {
        pb.authStore.clear()
        currentUser = nil
    }

    var isAuthenticated: Bool { currentUser != nil }
}
