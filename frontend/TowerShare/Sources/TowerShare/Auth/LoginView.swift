import SwiftUI

struct LoginView: View {
    @StateObject private var auth = AuthService.shared
    @State private var email = ""
    @State private var password = ""
    @State private var isSignUp = false
    @State private var displayName = ""

    var body: some View {
        VStack(spacing: 20) {
            Text("TowerShare")
                .font(.largeTitle.bold())

            if isSignUp {
                TextField("Display name", text: $displayName)
                    .textFieldStyle(.roundedBorder)
            }

            TextField("Email", text: $email)
                .textFieldStyle(.roundedBorder)
                .textInputAutocapitalization(.never)
                .keyboardType(.emailAddress)

            SecureField("Password", text: $password)
                .textFieldStyle(.roundedBorder)

            if let err = auth.error {
                Text(err).foregroundStyle(.red).font(.caption)
            }

            Button(isSignUp ? "Create account" : "Sign in") {
                Task {
                    if isSignUp {
                        await auth.signUp(email: email, password: password, displayName: displayName)
                    } else {
                        await auth.signIn(email: email, password: password)
                    }
                }
            }
            .buttonStyle(.borderedProminent)
            .disabled(auth.isLoading)

            Button(isSignUp ? "Already have an account? Sign in" : "No account? Sign up") {
                isSignUp.toggle()
            }
            .font(.footnote)
        }
        .padding(32)
    }
}
