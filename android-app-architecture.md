# WizardCore Android App - Local-First Architecture

## Project Structure

```
wizardcore-android/
├── app/
│   ├── src/
│   │   ├── main/
│   │   │   ├── java/com/wizardcore/
│   │   │   │   ├── WizardCoreApp.kt              # Application class
│   │   │   │   ├── di/                           # Dependency Injection
│   │   │   │   │   ├── AppModule.kt
│   │   │   │   │   ├── DatabaseModule.kt
│   │   │   │   │   ├── NetworkModule.kt
│   │   │   │   │   └── AnalyticsModule.kt
│   │   │   │   ├── data/
│   │   │   │   │   ├── local/                    # Local-first data layer
│   │   │   │   │   │   ├── database/
│   │   │   │   │   │   │   ├── AppDatabase.kt
│   │   │   │   │   │   │   ├── entities/
│   │   │   │   │   │   │   ├── daos/
│   │   │   │   │   │   │   └── migrations/
│   │   │   │   │   │   ├── datastore/
│   │   │   │   │   │   │   └── PreferencesDataStore.kt
│   │   │   │   │   │   └── repository/
│   │   │   │   │   │       ├── LocalRepository.kt
│   │   │   │   │   │       └── SyncRepository.kt
│   │   │   │   │   ├── remote/                   # Remote data layer
│   │   │   │   │   │   ├── api/
│   │   │   │   │   │   │   ├── WizardCoreApi.kt
│   │   │   │   │   │   │   ├── AuthApi.kt
│   │   │   │   │   │   │   └── Judge0Api.kt
│   │   │   │   │   │   └── repository/
│   │   │   │   │   │       └── RemoteRepository.kt
│   │   │   │   │   └── sync/                     # Sync engine
│   │   │   │   │       ├── SyncManager.kt
│   │   │   │   │       ├── SyncWorker.kt
│   │   │   │   │       └── ConflictResolver.kt
│   │   │   │   ├── domain/                       # Business logic
│   │   │   │   │   ├── model/
│   │   │   │   │   ├── repository/
│   │   │   │   │   └── usecase/
│   │   │   │   ├── presentation/
│   │   │   │   │   ├── ui/
│   │   │   │   │   │   ├── theme/
│   │   │   │   │   │   ├── components/
│   │   │   │   │   │   └── screens/
│   │   │   │   │   └── viewmodel/
│   │   │   │   └── analytics/                    # Analytics layer
│   │   │   │       ├── PostHogAnalytics.kt
│   │   │   │       ├── AnalyticsEvent.kt
│   │   │   │       └── AnalyticsWorker.kt
│   │   │   └── resources/
│   │   └── debug/                                # Debug builds
│   │       └── java/com/wizardcore/
│   │           └── debug/
│   │               └── DebugDataGenerator.kt
│   └── build.gradle.kts
├── build.gradle.kts
├── gradle.properties
├── local.properties
└── settings.gradle.kts
```

## Local-First Principles

1. **Data Always Available**: App works fully offline
2. **Background Sync**: Automatic sync when online
3. **Conflict Resolution**: Smart merging of changes
4. **Progressive Enhancement**: Better experience when online
5. **Privacy First**: Analytics opt-in, local data control

## Key Technologies

- **Kotlin** + **Coroutines** + **Flow**
- **Jetpack Compose** for UI
- **Room** for local database
- **DataStore** for preferences
- **WorkManager** for background tasks
- **Retrofit** + **OkHttp** for networking
- **Hilt** for dependency injection
- **PostHog** for analytics