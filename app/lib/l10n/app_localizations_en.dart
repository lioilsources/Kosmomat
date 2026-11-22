// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'KosmoMat';

  @override
  String get navToday => 'Today';

  @override
  String get navQuestions => 'Questions';

  @override
  String get navDisputes => 'Disputes';

  @override
  String get navHistory => 'History';

  @override
  String get navSearch => 'Search';

  @override
  String get mine => 'Mine';

  @override
  String get settings => 'Settings';

  @override
  String get todayChanged => 'What changed';

  @override
  String get todayCase => 'Case of the day';

  @override
  String get todayQuestion => 'Question of the day';

  @override
  String get todayDispute => 'Dispute of the day';

  @override
  String get searchHint => 'Search a term, theory, body…';

  @override
  String get searchEmpty => 'Nothing found. Try another word.';

  @override
  String get searchPrompt => 'Search works offline.';

  @override
  String get filterAll => 'All';

  @override
  String get level_kid => 'Child';

  @override
  String get level_lay => 'Lay';

  @override
  String get level_student => 'Student';

  @override
  String get lineage => 'Lineage';

  @override
  String get lineageOrigin => 'Origin';

  @override
  String get lineageAncestors => 'Where it comes from';

  @override
  String get lineageDescendants => 'What came of it';

  @override
  String get lineageCompetitors => 'Competitors';

  @override
  String get lineageThis => 'This card';

  @override
  String get opposedBy => 'Who disagrees and what contradicts it';

  @override
  String get supportedBy => 'What supports it';

  @override
  String get predictionBook => 'Prediction book';

  @override
  String get predMade => 'Made';

  @override
  String get predTested => 'Tested';

  @override
  String get predBefore => 'before measurement';

  @override
  String get predAfter => 'after measurement';

  @override
  String get predRetrodiction => 'retrodiction';

  @override
  String get outcome_confirmed => 'confirmed';

  @override
  String get outcome_refuted => 'refuted';

  @override
  String get outcome_partial => 'partial';

  @override
  String get outcome_open => 'open';

  @override
  String get outcome_untestable => 'untestable';

  @override
  String get sources => 'Sources';

  @override
  String get cardHistory => 'Card history';

  @override
  String lastVerified(String date) {
    return 'Last verified $date';
  }

  @override
  String get awaitsReview => 'Awaiting review';

  @override
  String get awaitsReviewDraft => 'Draft — not yet checked by a human.';

  @override
  String get notWrittenYet => 'Card not written yet';

  @override
  String get sourceTier => 'Source tier';

  @override
  String get sourceLocator => 'Location in the work';

  @override
  String get sourceCountry => 'Affiliation country';

  @override
  String get sourceOpen => 'Open source';

  @override
  String get sourceUnreviewed => 'Not peer reviewed';

  @override
  String get sourceRetracted => 'Retracted';

  @override
  String get sourceUnverified => 'Source metadata not yet verified';

  @override
  String get tier_T1 => 'T1 · peer-reviewed article';

  @override
  String get tier_T2 => 'T2 · preprint';

  @override
  String get tier_T3 => 'T3 · thesis';

  @override
  String get tier_T4 => 'T4 · agency, press release';

  @override
  String get tier_T5 => 'T5 · book';

  @override
  String get tier_T6 => 'T6 · popular (reception only)';

  @override
  String get status_established => 'established';

  @override
  String get status_strong => 'strong';

  @override
  String get status_contested => 'contested';

  @override
  String get status_speculative => 'speculative';

  @override
  String get status_refuted => 'refuted';

  @override
  String get status_historical => 'historical';

  @override
  String get status_philosophical => 'philosophical';

  @override
  String get statusHelp_established =>
      'At least two independent methods or teams, peer-reviewed sources and no live dispute in the peer-reviewed literature in the last 10 years.';

  @override
  String get statusHelp_strong =>
      'Rests on peer-reviewed sources but does not meet every rule for “established”.';

  @override
  String get statusHelp_contested =>
      'Peer-reviewed sources exist on both sides.';

  @override
  String get statusHelp_speculative => 'No decisive test exists yet.';

  @override
  String get statusHelp_refuted =>
      'A clearly stated prediction failed, or a directly contradicting observation has a peer-reviewed source.';

  @override
  String get statusHelp_historical =>
      'A historical model, described in its own time.';

  @override
  String get statusHelp_philosophical =>
      'Not decidable empirically; the card attributes positions to named thinkers.';

  @override
  String get lit_mainstream => 'mainstream';

  @override
  String get lit_minority => 'minority literature';

  @override
  String get lit_outside => 'outside peer review';

  @override
  String get litHelp =>
      'Where the debate takes place. “Minority” does not mean “wrong” and “mainstream” does not mean “certain” — strength of evidence is a separate label.';

  @override
  String get statusAxis => 'Strength of evidence';

  @override
  String get litAxis => 'Where it is debated';

  @override
  String get type_question => 'Question';

  @override
  String get type_observation => 'Observation';

  @override
  String get type_claim => 'Claim';

  @override
  String get type_theory => 'Theory';

  @override
  String get type_prediction => 'Prediction';

  @override
  String get type_case => 'Case study';

  @override
  String get type_person => 'Person';

  @override
  String get type_mission => 'Mission';

  @override
  String get type_term => 'Term';

  @override
  String get type_experiment => 'Try it yourself';

  @override
  String get type_era => 'Era';

  @override
  String get type_worldview => 'Worldview';

  @override
  String get type_body => 'Body';

  @override
  String get type_constellation => 'Constellation';

  @override
  String get listDisputesIntro =>
      'Theories by the same yardstick: claims, predictions, test results and open problems.';

  @override
  String get listHistoryIntro =>
      'How humanity changed its picture of the universe — including heresies that turned out wrong.';

  @override
  String get listQuestionsIntro =>
      'Big questions: myth, historical models and today\'s science kept apart.';

  @override
  String get listEmpty => 'Cards of this kind will arrive in later releases.';

  @override
  String get settingsLanguage => 'Content language';

  @override
  String get settingsLanguageSystem => 'System';

  @override
  String get settingsLevel => 'Default level';

  @override
  String get settingsEditorial => 'Editorial constitution';

  @override
  String get settingsCredits => 'Credits and licences';

  @override
  String get settingsPack => 'Content version';

  @override
  String get creditsWorks => 'Cited works';

  @override
  String get creditsMedia => 'Images';

  @override
  String get creditsMediaNone => 'This content version has no images yet.';

  @override
  String get creditsSoftware => 'Software licences';

  @override
  String get creditsEphemeris =>
      'Positions are computed with the astronomia library (MIT).';

  @override
  String get bookmark => 'Save';

  @override
  String get bookmarked => 'Saved';

  @override
  String get mineBookmarks => 'Bookmarks';

  @override
  String get mineEmpty =>
      'Nothing yet. Save cards with the star. Everything stays on your phone.';

  @override
  String unknownNode(String id) {
    return 'Card “$id” not found, try search.';
  }

  @override
  String packError(String error) {
    return 'Could not open content: $error';
  }
}
