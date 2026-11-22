// Maps pack enums to localized labels.
import '../l10n/app_localizations.dart';

extension Labels on AppLocalizations {
  String type(String t) => switch (t) {
        'question' => type_question,
        'observation' => type_observation,
        'claim' => type_claim,
        'theory' => type_theory,
        'prediction' => type_prediction,
        'case' => type_case,
        'person' => type_person,
        'mission' => type_mission,
        'term' => type_term,
        'experiment' => type_experiment,
        'era' => type_era,
        'worldview' => type_worldview,
        'body' => type_body,
        'constellation' => type_constellation,
        _ => t,
      };

  String status(String s) => switch (s) {
        'established' => status_established,
        'strong' => status_strong,
        'contested' => status_contested,
        'speculative' => status_speculative,
        'refuted' => status_refuted,
        'historical' => status_historical,
        'philosophical' => status_philosophical,
        _ => s,
      };

  String statusHelp(String s) => switch (s) {
        'established' => statusHelp_established,
        'strong' => statusHelp_strong,
        'contested' => statusHelp_contested,
        'speculative' => statusHelp_speculative,
        'refuted' => statusHelp_refuted,
        'historical' => statusHelp_historical,
        'philosophical' => statusHelp_philosophical,
        _ => s,
      };

  String literature(String l) => switch (l) {
        'mainstream' => lit_mainstream,
        'minority' => lit_minority,
        'outside' => lit_outside,
        _ => l,
      };

  String outcome(String o) => switch (o) {
        'confirmed' => outcome_confirmed,
        'refuted' => outcome_refuted,
        'partial' => outcome_partial,
        'open' => outcome_open,
        'untestable' => outcome_untestable,
        _ => o,
      };

  String tier(String t) => switch (t) {
        'T1' => tier_T1,
        'T2' => tier_T2,
        'T3' => tier_T3,
        'T4' => tier_T4,
        'T5' => tier_T5,
        'T6' => tier_T6,
        _ => t,
      };

  String level(String l) => switch (l) {
        'kid' => level_kid,
        'student' => level_student,
        _ => level_lay,
      };
}
