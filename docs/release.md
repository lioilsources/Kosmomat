# Release — ruční kroky

Pipeline z Distribution je nainstalovaná (`.github/workflows/ci.yml`, `release-ios.yml`,
`release-android.yml`, `app/ios/ExportOptions.plist`). Tyhle kroky skript neudělá:

1. **GitHub repo** `lioilsources/KosmoMAT` (public — kvůli GitHub Pages) a první push.
2. **Apple Developer:** App ID `com.ol1n.kosmomat` (Associated Domains až v M5 kvůli Universal Links),
   distribuční provisioning profil uložit do složky appky `Apple/KosmoMat/`.
3. **App Store Connect:** nová appka `KosmoMat`, bundle `com.ol1n.kosmomat`, kategorie Vzdělávání / Reference.
   Privacy Policy URL `https://lioilsources.github.io/KosmoMAT/privacy.html`,
   Support URL `https://lioilsources.github.io/KosmoMAT/support.html`, marketing `https://olin.now/kosmomat`.
4. **Firebase:** Android appka `com.ol1n.kosmomat` pro App Tester.
5. **Secrets:** `../Distribution/scripts/setup-gh-secrets.sh --repo lioilsources/KosmoMAT --app KosmoMat`.
6. **Android podpis:** do `app/android/app/build.gradle.kts` přidat načtení `key.properties`
   a release `signingConfig` (vzor Kiran) — zatím neuděláno.
7. **GitHub Pages:** `gh api -X POST repos/lioilsources/KosmoMAT/pages -f 'source[branch]=main' -f 'source[path]=/docs'`.
8. **Ikona** (FLUX přes SPARK) a stránka na olin.now (`apps/kosmomat/`) — viz skill `setup-flutter-app`.

Před releasem: externí odborná kontrola (plán §9) a všechny karty `verified_by: human`.
