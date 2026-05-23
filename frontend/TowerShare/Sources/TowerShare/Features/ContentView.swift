import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var auth: AuthService

    var body: some View {
        NavigationStack {
            Text("Welcome to TowerShare")
                .navigationTitle("TowerShare")
                .toolbar {
                    ToolbarItem(placement: .topBarTrailing) {
                        Button("Sign out") { auth.signOut() }
                    }
                }
        }
    }
}
