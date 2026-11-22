// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Czech (`cs`).
class AppLocalizationsCs extends AppLocalizations {
  AppLocalizationsCs([String locale = 'cs']) : super(locale);

  @override
  String get appTitle => 'KosmoMat';

  @override
  String get navToday => 'Dnes';

  @override
  String get navQuestions => 'Otázky';

  @override
  String get navDisputes => 'Spory';

  @override
  String get navHistory => 'Dějiny';

  @override
  String get navSearch => 'Hledat';

  @override
  String get mine => 'Moje';

  @override
  String get settings => 'Nastavení';

  @override
  String get todayChanged => 'Co se změnilo';

  @override
  String get todayCase => 'Případ dne';

  @override
  String get todayQuestion => 'Otázka dne';

  @override
  String get todayDispute => 'Spor dne';

  @override
  String get searchHint => 'Hledej pojem, teorii, těleso…';

  @override
  String get searchEmpty => 'Nic jsme nenašli. Zkus jiné slovo.';

  @override
  String get searchPrompt => 'Hledání funguje i bez připojení.';

  @override
  String get filterAll => 'Vše';

  @override
  String get level_kid => 'Dítě';

  @override
  String get level_lay => 'Laik';

  @override
  String get level_student => 'Student';

  @override
  String get lineage => 'Rodokmen';

  @override
  String get lineageOrigin => 'Počátek';

  @override
  String get lineageAncestors => 'Z čeho vychází';

  @override
  String get lineageDescendants => 'Co z toho vzešlo';

  @override
  String get lineageCompetitors => 'Konkurenti';

  @override
  String get lineageThis => 'Tato karta';

  @override
  String get opposedBy => 'Kdo oponuje a co tomu odporuje';

  @override
  String get supportedBy => 'O co se opírá';

  @override
  String get predictionBook => 'Kniha předpovědí';

  @override
  String get predMade => 'Vysloveno';

  @override
  String get predTested => 'Testováno';

  @override
  String get predBefore => 'před měřením';

  @override
  String get predAfter => 'po měření';

  @override
  String get predRetrodiction => 'retrodikce';

  @override
  String get outcome_confirmed => 'potvrzeno';

  @override
  String get outcome_refuted => 'vyvráceno';

  @override
  String get outcome_partial => 'částečně';

  @override
  String get outcome_open => 'otevřené';

  @override
  String get outcome_untestable => 'netestovatelné';

  @override
  String get sources => 'Zdroje';

  @override
  String get cardHistory => 'Historie karty';

  @override
  String lastVerified(String date) {
    return 'Naposledy ověřeno $date';
  }

  @override
  String get awaitsReview => 'Čeká na revizi';

  @override
  String get awaitsReviewDraft => 'Návrh — kartu ještě nezkontroloval člověk.';

  @override
  String get notWrittenYet => 'Karta zatím nevznikla';

  @override
  String get sourceTier => 'Úroveň zdroje';

  @override
  String get sourceLocator => 'Místo v díle';

  @override
  String get sourceCountry => 'Země pracoviště';

  @override
  String get sourceOpen => 'Otevřít zdroj';

  @override
  String get sourceUnreviewed => 'Nerecenzováno';

  @override
  String get sourceRetracted => 'Staženo';

  @override
  String get sourceUnverified => 'Metadata zdroje zatím neověřena';

  @override
  String get tier_T1 => 'T1 · recenzovaný článek';

  @override
  String get tier_T2 => 'T2 · preprint';

  @override
  String get tier_T3 => 'T3 · dizertace';

  @override
  String get tier_T4 => 'T4 · agentura, tisková zpráva';

  @override
  String get tier_T5 => 'T5 · kniha';

  @override
  String get tier_T6 => 'T6 · popularizace (jen ohlas)';

  @override
  String get status_established => 'ověřené';

  @override
  String get status_strong => 'silné';

  @override
  String get status_contested => 'sporné';

  @override
  String get status_speculative => 'spekulativní';

  @override
  String get status_refuted => 'vyvrácené';

  @override
  String get status_historical => 'historické';

  @override
  String get status_philosophical => 'filozofické';

  @override
  String get statusHelp_established =>
      'Aspoň dvě nezávislé metody nebo týmy, recenzované zdroje a žádný živý spor v recenzované literatuře za posledních 10 let.';

  @override
  String get statusHelp_strong =>
      'Opírá se o recenzované zdroje, ale nesplňuje všechna pravidla pro „ověřené“.';

  @override
  String get statusHelp_contested =>
      'Recenzované zdroje existují na obou stranách sporu.';

  @override
  String get statusHelp_speculative => 'Zatím neexistuje rozhodující test.';

  @override
  String get statusHelp_refuted =>
      'Selhala jasně formulovaná předpověď, nebo existuje přímo odporující pozorování s recenzovaným zdrojem.';

  @override
  String get statusHelp_historical =>
      'Historický model; popisujeme ho v jeho době.';

  @override
  String get statusHelp_philosophical =>
      'Empiricky nerozhodnutelné; karta přiřazuje postoje konkrétním myslitelům.';

  @override
  String get lit_mainstream => 'hlavní proud';

  @override
  String get lit_minority => 'menšinová literatura';

  @override
  String get lit_outside => 'mimo recenzovanou literaturu';

  @override
  String get litHelp =>
      'Kde se o tom vede debata. „Menšinový“ neznamená „špatný“ a „hlavní proud“ neznamená „jistý“ — síla důkazů je samostatná nálepka.';

  @override
  String get statusAxis => 'Síla důkazů';

  @override
  String get litAxis => 'Kde se o tom debatuje';

  @override
  String get type_question => 'Otázka';

  @override
  String get type_observation => 'Pozorování';

  @override
  String get type_claim => 'Tvrzení';

  @override
  String get type_theory => 'Teorie';

  @override
  String get type_prediction => 'Předpověď';

  @override
  String get type_case => 'Případová studie';

  @override
  String get type_person => 'Osobnost';

  @override
  String get type_mission => 'Mise';

  @override
  String get type_term => 'Pojem';

  @override
  String get type_experiment => 'Ověř si sám';

  @override
  String get type_era => 'Epocha';

  @override
  String get type_worldview => 'Pohled na svět';

  @override
  String get type_body => 'Těleso';

  @override
  String get type_constellation => 'Souhvězdí';

  @override
  String get listDisputesIntro =>
      'Teorie stejným metrem: tvrzení, předpovědi, výsledky testů a otevřené problémy.';

  @override
  String get listHistoryIntro =>
      'Jak lidstvo měnilo obraz vesmíru — včetně omylů, které se nepotvrdily.';

  @override
  String get listQuestionsIntro =>
      'Velké otázky: mýtus, historické modely a dnešní věda zvlášť.';

  @override
  String get listEmpty => 'Karty tohoto typu přibudou v dalších vydáních.';

  @override
  String get settingsLanguage => 'Jazyk obsahu';

  @override
  String get settingsLanguageSystem => 'Podle systému';

  @override
  String get settingsLevel => 'Výchozí úroveň';

  @override
  String get settingsEditorial => 'Redakční ústava';

  @override
  String get settingsCredits => 'Kredity a licence';

  @override
  String get settingsPack => 'Verze obsahu';

  @override
  String get creditsWorks => 'Citovaná díla';

  @override
  String get creditsMedia => 'Obrázky';

  @override
  String get creditsMediaNone => 'Tato verze obsahu zatím neobsahuje obrázky.';

  @override
  String get creditsSoftware => 'Licence softwaru';

  @override
  String get creditsEphemeris =>
      'Polohy těles počítá knihovna astronomia (MIT).';

  @override
  String get bookmark => 'Uložit';

  @override
  String get bookmarked => 'Uloženo';

  @override
  String get mineBookmarks => 'Záložky';

  @override
  String get mineEmpty =>
      'Zatím nic. Karty si uložíš hvězdičkou. Všechno zůstává jen v telefonu.';

  @override
  String unknownNode(String id) {
    return 'Kartu „$id“ jsme nenašli, zkus hledání.';
  }

  @override
  String packError(String error) {
    return 'Obsah se nepodařilo otevřít: $error';
  }
}
