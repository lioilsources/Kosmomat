import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_cs.dart';
import 'app_localizations_en.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('cs'),
    Locale('en'),
  ];

  /// No description provided for @appTitle.
  ///
  /// In cs, this message translates to:
  /// **'KosmoMat'**
  String get appTitle;

  /// No description provided for @navToday.
  ///
  /// In cs, this message translates to:
  /// **'Dnes'**
  String get navToday;

  /// No description provided for @navQuestions.
  ///
  /// In cs, this message translates to:
  /// **'Otázky'**
  String get navQuestions;

  /// No description provided for @navDisputes.
  ///
  /// In cs, this message translates to:
  /// **'Spory'**
  String get navDisputes;

  /// No description provided for @navHistory.
  ///
  /// In cs, this message translates to:
  /// **'Dějiny'**
  String get navHistory;

  /// No description provided for @navSearch.
  ///
  /// In cs, this message translates to:
  /// **'Hledat'**
  String get navSearch;

  /// No description provided for @mine.
  ///
  /// In cs, this message translates to:
  /// **'Moje'**
  String get mine;

  /// No description provided for @settings.
  ///
  /// In cs, this message translates to:
  /// **'Nastavení'**
  String get settings;

  /// No description provided for @todayChanged.
  ///
  /// In cs, this message translates to:
  /// **'Co se změnilo'**
  String get todayChanged;

  /// No description provided for @todayCase.
  ///
  /// In cs, this message translates to:
  /// **'Případ dne'**
  String get todayCase;

  /// No description provided for @todayQuestion.
  ///
  /// In cs, this message translates to:
  /// **'Otázka dne'**
  String get todayQuestion;

  /// No description provided for @todayDispute.
  ///
  /// In cs, this message translates to:
  /// **'Spor dne'**
  String get todayDispute;

  /// No description provided for @searchHint.
  ///
  /// In cs, this message translates to:
  /// **'Hledej pojem, teorii, těleso…'**
  String get searchHint;

  /// No description provided for @searchEmpty.
  ///
  /// In cs, this message translates to:
  /// **'Nic jsme nenašli. Zkus jiné slovo.'**
  String get searchEmpty;

  /// No description provided for @searchPrompt.
  ///
  /// In cs, this message translates to:
  /// **'Hledání funguje i bez připojení.'**
  String get searchPrompt;

  /// No description provided for @filterAll.
  ///
  /// In cs, this message translates to:
  /// **'Vše'**
  String get filterAll;

  /// No description provided for @level_kid.
  ///
  /// In cs, this message translates to:
  /// **'Dítě'**
  String get level_kid;

  /// No description provided for @level_lay.
  ///
  /// In cs, this message translates to:
  /// **'Laik'**
  String get level_lay;

  /// No description provided for @level_student.
  ///
  /// In cs, this message translates to:
  /// **'Student'**
  String get level_student;

  /// No description provided for @lineage.
  ///
  /// In cs, this message translates to:
  /// **'Rodokmen'**
  String get lineage;

  /// No description provided for @lineageOrigin.
  ///
  /// In cs, this message translates to:
  /// **'Počátek'**
  String get lineageOrigin;

  /// No description provided for @lineageAncestors.
  ///
  /// In cs, this message translates to:
  /// **'Z čeho vychází'**
  String get lineageAncestors;

  /// No description provided for @lineageDescendants.
  ///
  /// In cs, this message translates to:
  /// **'Co z toho vzešlo'**
  String get lineageDescendants;

  /// No description provided for @lineageCompetitors.
  ///
  /// In cs, this message translates to:
  /// **'Konkurenti'**
  String get lineageCompetitors;

  /// No description provided for @lineageThis.
  ///
  /// In cs, this message translates to:
  /// **'Tato karta'**
  String get lineageThis;

  /// No description provided for @opposedBy.
  ///
  /// In cs, this message translates to:
  /// **'Kdo oponuje a co tomu odporuje'**
  String get opposedBy;

  /// No description provided for @supportedBy.
  ///
  /// In cs, this message translates to:
  /// **'O co se opírá'**
  String get supportedBy;

  /// No description provided for @predictionBook.
  ///
  /// In cs, this message translates to:
  /// **'Kniha předpovědí'**
  String get predictionBook;

  /// No description provided for @predMade.
  ///
  /// In cs, this message translates to:
  /// **'Vysloveno'**
  String get predMade;

  /// No description provided for @predTested.
  ///
  /// In cs, this message translates to:
  /// **'Testováno'**
  String get predTested;

  /// No description provided for @predBefore.
  ///
  /// In cs, this message translates to:
  /// **'před měřením'**
  String get predBefore;

  /// No description provided for @predAfter.
  ///
  /// In cs, this message translates to:
  /// **'po měření'**
  String get predAfter;

  /// No description provided for @predRetrodiction.
  ///
  /// In cs, this message translates to:
  /// **'retrodikce'**
  String get predRetrodiction;

  /// No description provided for @outcome_confirmed.
  ///
  /// In cs, this message translates to:
  /// **'potvrzeno'**
  String get outcome_confirmed;

  /// No description provided for @outcome_refuted.
  ///
  /// In cs, this message translates to:
  /// **'vyvráceno'**
  String get outcome_refuted;

  /// No description provided for @outcome_partial.
  ///
  /// In cs, this message translates to:
  /// **'částečně'**
  String get outcome_partial;

  /// No description provided for @outcome_open.
  ///
  /// In cs, this message translates to:
  /// **'otevřené'**
  String get outcome_open;

  /// No description provided for @outcome_untestable.
  ///
  /// In cs, this message translates to:
  /// **'netestovatelné'**
  String get outcome_untestable;

  /// No description provided for @sources.
  ///
  /// In cs, this message translates to:
  /// **'Zdroje'**
  String get sources;

  /// No description provided for @cardHistory.
  ///
  /// In cs, this message translates to:
  /// **'Historie karty'**
  String get cardHistory;

  /// No description provided for @lastVerified.
  ///
  /// In cs, this message translates to:
  /// **'Naposledy ověřeno {date}'**
  String lastVerified(String date);

  /// No description provided for @awaitsReview.
  ///
  /// In cs, this message translates to:
  /// **'Čeká na revizi'**
  String get awaitsReview;

  /// No description provided for @awaitsReviewDraft.
  ///
  /// In cs, this message translates to:
  /// **'Návrh — kartu ještě nezkontroloval člověk.'**
  String get awaitsReviewDraft;

  /// No description provided for @notWrittenYet.
  ///
  /// In cs, this message translates to:
  /// **'Karta zatím nevznikla'**
  String get notWrittenYet;

  /// No description provided for @sourceTier.
  ///
  /// In cs, this message translates to:
  /// **'Úroveň zdroje'**
  String get sourceTier;

  /// No description provided for @sourceLocator.
  ///
  /// In cs, this message translates to:
  /// **'Místo v díle'**
  String get sourceLocator;

  /// No description provided for @sourceCountry.
  ///
  /// In cs, this message translates to:
  /// **'Země pracoviště'**
  String get sourceCountry;

  /// No description provided for @sourceOpen.
  ///
  /// In cs, this message translates to:
  /// **'Otevřít zdroj'**
  String get sourceOpen;

  /// No description provided for @sourceUnreviewed.
  ///
  /// In cs, this message translates to:
  /// **'Nerecenzováno'**
  String get sourceUnreviewed;

  /// No description provided for @sourceRetracted.
  ///
  /// In cs, this message translates to:
  /// **'Staženo'**
  String get sourceRetracted;

  /// No description provided for @sourceUnverified.
  ///
  /// In cs, this message translates to:
  /// **'Metadata zdroje zatím neověřena'**
  String get sourceUnverified;

  /// No description provided for @tier_T1.
  ///
  /// In cs, this message translates to:
  /// **'T1 · recenzovaný článek'**
  String get tier_T1;

  /// No description provided for @tier_T2.
  ///
  /// In cs, this message translates to:
  /// **'T2 · preprint'**
  String get tier_T2;

  /// No description provided for @tier_T3.
  ///
  /// In cs, this message translates to:
  /// **'T3 · dizertace'**
  String get tier_T3;

  /// No description provided for @tier_T4.
  ///
  /// In cs, this message translates to:
  /// **'T4 · agentura, tisková zpráva'**
  String get tier_T4;

  /// No description provided for @tier_T5.
  ///
  /// In cs, this message translates to:
  /// **'T5 · kniha'**
  String get tier_T5;

  /// No description provided for @tier_T6.
  ///
  /// In cs, this message translates to:
  /// **'T6 · popularizace (jen ohlas)'**
  String get tier_T6;

  /// No description provided for @status_established.
  ///
  /// In cs, this message translates to:
  /// **'ověřené'**
  String get status_established;

  /// No description provided for @status_strong.
  ///
  /// In cs, this message translates to:
  /// **'silné'**
  String get status_strong;

  /// No description provided for @status_contested.
  ///
  /// In cs, this message translates to:
  /// **'sporné'**
  String get status_contested;

  /// No description provided for @status_speculative.
  ///
  /// In cs, this message translates to:
  /// **'spekulativní'**
  String get status_speculative;

  /// No description provided for @status_refuted.
  ///
  /// In cs, this message translates to:
  /// **'vyvrácené'**
  String get status_refuted;

  /// No description provided for @status_historical.
  ///
  /// In cs, this message translates to:
  /// **'historické'**
  String get status_historical;

  /// No description provided for @status_philosophical.
  ///
  /// In cs, this message translates to:
  /// **'filozofické'**
  String get status_philosophical;

  /// No description provided for @statusHelp_established.
  ///
  /// In cs, this message translates to:
  /// **'Aspoň dvě nezávislé metody nebo týmy, recenzované zdroje a žádný živý spor v recenzované literatuře za posledních 10 let.'**
  String get statusHelp_established;

  /// No description provided for @statusHelp_strong.
  ///
  /// In cs, this message translates to:
  /// **'Opírá se o recenzované zdroje, ale nesplňuje všechna pravidla pro „ověřené“.'**
  String get statusHelp_strong;

  /// No description provided for @statusHelp_contested.
  ///
  /// In cs, this message translates to:
  /// **'Recenzované zdroje existují na obou stranách sporu.'**
  String get statusHelp_contested;

  /// No description provided for @statusHelp_speculative.
  ///
  /// In cs, this message translates to:
  /// **'Zatím neexistuje rozhodující test.'**
  String get statusHelp_speculative;

  /// No description provided for @statusHelp_refuted.
  ///
  /// In cs, this message translates to:
  /// **'Selhala jasně formulovaná předpověď, nebo existuje přímo odporující pozorování s recenzovaným zdrojem.'**
  String get statusHelp_refuted;

  /// No description provided for @statusHelp_historical.
  ///
  /// In cs, this message translates to:
  /// **'Historický model; popisujeme ho v jeho době.'**
  String get statusHelp_historical;

  /// No description provided for @statusHelp_philosophical.
  ///
  /// In cs, this message translates to:
  /// **'Empiricky nerozhodnutelné; karta přiřazuje postoje konkrétním myslitelům.'**
  String get statusHelp_philosophical;

  /// No description provided for @lit_mainstream.
  ///
  /// In cs, this message translates to:
  /// **'hlavní proud'**
  String get lit_mainstream;

  /// No description provided for @lit_minority.
  ///
  /// In cs, this message translates to:
  /// **'menšinová literatura'**
  String get lit_minority;

  /// No description provided for @lit_outside.
  ///
  /// In cs, this message translates to:
  /// **'mimo recenzovanou literaturu'**
  String get lit_outside;

  /// No description provided for @litHelp.
  ///
  /// In cs, this message translates to:
  /// **'Kde se o tom vede debata. „Menšinový“ neznamená „špatný“ a „hlavní proud“ neznamená „jistý“ — síla důkazů je samostatná nálepka.'**
  String get litHelp;

  /// No description provided for @statusAxis.
  ///
  /// In cs, this message translates to:
  /// **'Síla důkazů'**
  String get statusAxis;

  /// No description provided for @litAxis.
  ///
  /// In cs, this message translates to:
  /// **'Kde se o tom debatuje'**
  String get litAxis;

  /// No description provided for @type_question.
  ///
  /// In cs, this message translates to:
  /// **'Otázka'**
  String get type_question;

  /// No description provided for @type_observation.
  ///
  /// In cs, this message translates to:
  /// **'Pozorování'**
  String get type_observation;

  /// No description provided for @type_claim.
  ///
  /// In cs, this message translates to:
  /// **'Tvrzení'**
  String get type_claim;

  /// No description provided for @type_theory.
  ///
  /// In cs, this message translates to:
  /// **'Teorie'**
  String get type_theory;

  /// No description provided for @type_prediction.
  ///
  /// In cs, this message translates to:
  /// **'Předpověď'**
  String get type_prediction;

  /// No description provided for @type_case.
  ///
  /// In cs, this message translates to:
  /// **'Případová studie'**
  String get type_case;

  /// No description provided for @type_person.
  ///
  /// In cs, this message translates to:
  /// **'Osobnost'**
  String get type_person;

  /// No description provided for @type_mission.
  ///
  /// In cs, this message translates to:
  /// **'Mise'**
  String get type_mission;

  /// No description provided for @type_term.
  ///
  /// In cs, this message translates to:
  /// **'Pojem'**
  String get type_term;

  /// No description provided for @type_experiment.
  ///
  /// In cs, this message translates to:
  /// **'Ověř si sám'**
  String get type_experiment;

  /// No description provided for @type_era.
  ///
  /// In cs, this message translates to:
  /// **'Epocha'**
  String get type_era;

  /// No description provided for @type_worldview.
  ///
  /// In cs, this message translates to:
  /// **'Pohled na svět'**
  String get type_worldview;

  /// No description provided for @type_body.
  ///
  /// In cs, this message translates to:
  /// **'Těleso'**
  String get type_body;

  /// No description provided for @type_constellation.
  ///
  /// In cs, this message translates to:
  /// **'Souhvězdí'**
  String get type_constellation;

  /// No description provided for @listDisputesIntro.
  ///
  /// In cs, this message translates to:
  /// **'Teorie stejným metrem: tvrzení, předpovědi, výsledky testů a otevřené problémy.'**
  String get listDisputesIntro;

  /// No description provided for @listHistoryIntro.
  ///
  /// In cs, this message translates to:
  /// **'Jak lidstvo měnilo obraz vesmíru — včetně omylů, které se nepotvrdily.'**
  String get listHistoryIntro;

  /// No description provided for @listQuestionsIntro.
  ///
  /// In cs, this message translates to:
  /// **'Velké otázky: mýtus, historické modely a dnešní věda zvlášť.'**
  String get listQuestionsIntro;

  /// No description provided for @listEmpty.
  ///
  /// In cs, this message translates to:
  /// **'Karty tohoto typu přibudou v dalších vydáních.'**
  String get listEmpty;

  /// No description provided for @settingsLanguage.
  ///
  /// In cs, this message translates to:
  /// **'Jazyk obsahu'**
  String get settingsLanguage;

  /// No description provided for @settingsLanguageSystem.
  ///
  /// In cs, this message translates to:
  /// **'Podle systému'**
  String get settingsLanguageSystem;

  /// No description provided for @settingsLevel.
  ///
  /// In cs, this message translates to:
  /// **'Výchozí úroveň'**
  String get settingsLevel;

  /// No description provided for @settingsEditorial.
  ///
  /// In cs, this message translates to:
  /// **'Redakční ústava'**
  String get settingsEditorial;

  /// No description provided for @settingsCredits.
  ///
  /// In cs, this message translates to:
  /// **'Kredity a licence'**
  String get settingsCredits;

  /// No description provided for @settingsPack.
  ///
  /// In cs, this message translates to:
  /// **'Verze obsahu'**
  String get settingsPack;

  /// No description provided for @creditsWorks.
  ///
  /// In cs, this message translates to:
  /// **'Citovaná díla'**
  String get creditsWorks;

  /// No description provided for @creditsMedia.
  ///
  /// In cs, this message translates to:
  /// **'Obrázky'**
  String get creditsMedia;

  /// No description provided for @creditsMediaNone.
  ///
  /// In cs, this message translates to:
  /// **'Tato verze obsahu zatím neobsahuje obrázky.'**
  String get creditsMediaNone;

  /// No description provided for @creditsSoftware.
  ///
  /// In cs, this message translates to:
  /// **'Licence softwaru'**
  String get creditsSoftware;

  /// No description provided for @creditsEphemeris.
  ///
  /// In cs, this message translates to:
  /// **'Polohy těles počítá knihovna astronomia (MIT).'**
  String get creditsEphemeris;

  /// No description provided for @bookmark.
  ///
  /// In cs, this message translates to:
  /// **'Uložit'**
  String get bookmark;

  /// No description provided for @bookmarked.
  ///
  /// In cs, this message translates to:
  /// **'Uloženo'**
  String get bookmarked;

  /// No description provided for @mineBookmarks.
  ///
  /// In cs, this message translates to:
  /// **'Záložky'**
  String get mineBookmarks;

  /// No description provided for @mineEmpty.
  ///
  /// In cs, this message translates to:
  /// **'Zatím nic. Karty si uložíš hvězdičkou. Všechno zůstává jen v telefonu.'**
  String get mineEmpty;

  /// No description provided for @unknownNode.
  ///
  /// In cs, this message translates to:
  /// **'Kartu „{id}“ jsme nenašli, zkus hledání.'**
  String unknownNode(String id);

  /// No description provided for @packError.
  ///
  /// In cs, this message translates to:
  /// **'Obsah se nepodařilo otevřít: {error}'**
  String packError(String error);
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['cs', 'en'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'cs':
      return AppLocalizationsCs();
    case 'en':
      return AppLocalizationsEn();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
