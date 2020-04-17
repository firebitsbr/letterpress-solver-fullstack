package solver

import "encoding/json"

func UnmarshalDuolingoSession(data []byte) (DuolingoSession, error) {
	var r DuolingoSession
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DuolingoSession) MarshalDuolingoSession() ([]byte, error) {
	return json.Marshal(r)
}

type DuolingoSession struct {
	FromLanguage             FromLanguageEnum    `json:"fromLanguage"`
	LevelSessionIndex        int64               `json:"levelSessionIndex"`
	LearningLanguage         FromLanguageEnum    `json:"learningLanguage"`
	Explanations             Explanations        `json:"explanations"`
	TrackingProperties       TrackingProperties  `json:"trackingProperties"`
	LessonIndex              int64               `json:"lessonIndex"`
	Metadata                 SessionMetadata     `json:"metadata"`
	ProgressUpdates          []ProgressUpdate    `json:"progressUpdates"`
	Challenges               []Challenge         `json:"challenges"`
	LevelIndex               int64               `json:"levelIndex"`
	AdaptiveChallenges       []AdaptiveChallenge `json:"adaptiveChallenges"`
	SessionStartExperiments  []interface{}       `json:"sessionStartExperiments"`
	Type                     string              `json:"type"`
	ID                       string              `json:"id"`
	ChallengeTimeTakenCutoff int64               `json:"challengeTimeTakenCutoff"`
	SkillID                  string              `json:"skillId"`
}

type AdaptiveChallenge struct {
	SentenceDiscussionID         string                       `json:"sentenceDiscussionId"`
	SentenceID                   string                       `json:"sentenceId"`
	Prompt                       string                       `json:"prompt"`
	SourceLanguage               FromLanguageEnum             `json:"sourceLanguage"`
	TTS                          string                       `json:"tts"`
	WrongTokens                  []string                     `json:"wrongTokens"`
	CompactTranslations          []string                     `json:"compactTranslations"`
	ProgressUpdates              []interface{}                `json:"progressUpdates"`
	CorrectIndices               []int64                      `json:"correctIndices"`
	Tokens                       []AdaptiveChallengeToken     `json:"tokens"`
	Grader                       Grader                       `json:"grader"`
	CorrectTokens                []string                     `json:"correctTokens"`
	TargetLanguage               FromLanguageEnum             `json:"targetLanguage"`
	ID                           string                       `json:"id"`
	NewWords                     []string                     `json:"newWords"`
	ChallengeGeneratorIdentifier ChallengeGeneratorIdentifier `json:"challengeGeneratorIdentifier"`
	CorrectSolutions             []string                     `json:"correctSolutions"`
	Choices                      []Choice                     `json:"choices"`
	Type                         MetadataType                 `json:"type"`
	Metadata                     AdaptiveChallengeMetadata    `json:"metadata"`
}

type ChallengeGeneratorIdentifier struct {
	GeneratorID  string `json:"generatorId"`
	SpecificType string `json:"specificType"`
}

type Choice struct {
	TTS  *string `json:"tts,omitempty"`
	Text string  `json:"text"`
}

type Grader struct {
	WhitespaceDelimited bool             `json:"whitespaceDelimited"`
	Version             int64            `json:"version"`
	Vertices            [][]Vertex       `json:"vertices"`
	Language            FromLanguageEnum `json:"language"`
}

type Vertex struct {
	To      int64       `json:"to"`
	Lenient string      `json:"lenient"`
	Orig    *string     `json:"orig,omitempty"`
	Type    *VertexType `json:"type,omitempty"`
	Auto    *bool       `json:"auto,omitempty"`
}

type AdaptiveChallengeMetadata struct {
	TaggedKcIDS         []TaggedKcID         `json:"tagged_kc_ids"`
	Sentence            string               `json:"sentence"`
	SolutionKey         string               `json:"solution_key"`
	SourceLanguage      FromLanguageEnum     `json:"source_language"`
	ActivityUUID        string               `json:"activity_uuid"`
	WrongTokens         []string             `json:"wrong_tokens"`
	LexemeIDSToUpdate   []string             `json:"lexeme_ids_to_update"`
	UnknownWords        []string             `json:"unknown_words"`
	SpecificType        SpecificType         `json:"specific_type"`
	UUID                string               `json:"uuid"`
	WrongLexemes        []string             `json:"wrong_lexemes"`
	GeneratorID         string               `json:"generator_id"`
	Text                string               `json:"text"`
	Type                MetadataType         `json:"type"`
	NewExplanationIDS   []interface{}        `json:"new_explanation_ids"`
	LexemeIDS           []string             `json:"lexeme_ids"`
	TeachesLexemeIDS    []string             `json:"teaches_lexeme_ids"`
	Tokens              []string             `json:"tokens"`
	SentenceID          string               `json:"sentence_id"`
	DiscussionCount     int64                `json:"discussion_count"`
	TargetLanguageName  string               `json:"target_language_name"`
	Translation         string               `json:"translation"`
	FromLanguage        FromLanguageEnum     `json:"from_language"`
	HasTTS              bool                 `json:"has_tts"`
	LearningLanguage    LearningLanguageEnum `json:"learning_language"`
	HasAccents          bool                 `json:"has_accents"`
	Highlight           []string             `json:"highlight"`
	KnowledgeComponents []string             `json:"knowledge_components"`
	TargetLanguage      LearningLanguageEnum `json:"target_language"`
}

type TaggedKcID struct {
	KcTypeStr KcTypeStr `json:"kc_type_str"`
	LegacyID  string    `json:"legacy_id"`
}

type AdaptiveChallengeToken struct {
	HintTable *HintTable `json:"hintTable,omitempty"`
	Value     string     `json:"value"`
}

type HintTable struct {
	Headers []string `json:"headers"`
	Rows    [][]Row  `json:"rows"`
}

type Row struct {
	Colspan int64   `json:"colspan"`
	Hint    *string `json:"hint,omitempty"`
}

type Challenge struct {
	SentenceDiscussionID         string                       `json:"sentenceDiscussionId"`
	SentenceID                   string                       `json:"sentenceId"`
	Prompt                       string                       `json:"prompt"`
	SourceLanguage               *FromLanguageEnum            `json:"sourceLanguage,omitempty"`
	TTS                          string                       `json:"tts"`
	WrongTokens                  []string                     `json:"wrongTokens"`
	CompactTranslations          []string                     `json:"compactTranslations"`
	ProgressUpdates              []interface{}                `json:"progressUpdates"`
	CorrectIndices               []int64                      `json:"correctIndices"`
	Tokens                       []ChallengeToken             `json:"tokens"`
	Grader                       *Grader                      `json:"grader,omitempty"`
	CorrectTokens                []string                     `json:"correctTokens"`
	TargetLanguage               *FromLanguageEnum            `json:"targetLanguage,omitempty"`
	ID                           string                       `json:"id"`
	NewWords                     []string                     `json:"newWords"`
	ChallengeGeneratorIdentifier ChallengeGeneratorIdentifier `json:"challengeGeneratorIdentifier"`
	CorrectSolutions             []string                     `json:"correctSolutions"`
	Choices                      []Choice                     `json:"choices"`
	Type                         ChallengeType                `json:"type"`
	Metadata                     ChallengeMetadata            `json:"metadata"`
	SolutionTranslation          *string                      `json:"solutionTranslation,omitempty"`
	SlowTTS                      *string                      `json:"slowTts,omitempty"`
	SoundID                      *string                      `json:"soundId,omitempty"`
	Threshold                    *float64                     `json:"threshold,omitempty"`
}

type ChallengeMetadata struct {
	TaggedKcIDS         []TaggedKcID          `json:"tagged_kc_ids"`
	NonCharacterTTS     *NonCharacterTTS      `json:"non_character_tts,omitempty"`
	Sentence            *string               `json:"sentence,omitempty"`
	SolutionKey         string                `json:"solution_key"`
	SourceLanguage      LearningLanguageEnum  `json:"source_language"`
	ActivityUUID        string                `json:"activity_uuid"`
	WrongTokens         []string              `json:"wrong_tokens"`
	LexemeIDSToUpdate   []string              `json:"lexeme_ids_to_update"`
	UnknownWords        []string              `json:"unknown_words"`
	SpecificType        SpecificType          `json:"specific_type"`
	UUID                string                `json:"uuid"`
	WrongLexemes        []string              `json:"wrong_lexemes"`
	NumComments         int64                 `json:"num_comments"`
	GeneratorID         string                `json:"generator_id"`
	Text                string                `json:"text"`
	Type                MetadataType          `json:"type"`
	NewExplanationIDS   []interface{}         `json:"new_explanation_ids"`
	GenericLexemeMap    Explanations          `json:"generic_lexeme_map"`
	LexemeIDS           []string              `json:"lexeme_ids"`
	TeachesLexemeIDS    []string              `json:"teaches_lexeme_ids"`
	Tokens              []string              `json:"tokens"`
	SentenceID          string                `json:"sentence_id"`
	DiscussionCount     int64                 `json:"discussion_count"`
	TargetLanguageName  *TargetLanguageName   `json:"target_language_name,omitempty"`
	Highlight           []string              `json:"highlight"`
	Translation         *string               `json:"translation,omitempty"`
	FromLanguage        FromLanguageEnum      `json:"from_language"`
	HasTTS              *bool                 `json:"has_tts,omitempty"`
	LearningLanguage    LearningLanguageEnum  `json:"learning_language"`
	HasAccents          *bool                 `json:"has_accents,omitempty"`
	LexemesToUpdate     []string              `json:"lexemes_to_update"`
	KnowledgeComponents []string              `json:"knowledge_components"`
	TargetLanguage      *FromLanguageEnum     `json:"target_language,omitempty"`
	SolutionTranslation *string               `json:"solution_translation,omitempty"`
	LanguageName        *string               `json:"language_name,omitempty"`
	CorrectSolutions    *string               `json:"correct_solutions,omitempty"`
	Language            *LearningLanguageEnum `json:"language,omitempty"`
	Threshold           *float64              `json:"threshold,omitempty"`
	WsLanguage          *FromLanguageEnum     `json:"ws_language,omitempty"`
	StrippedText        *string               `json:"stripped_text,omitempty"`
	WsThreshold         *float64              `json:"ws_threshold,omitempty"`
	SoundID             *string               `json:"sound_id,omitempty"`
}

type Explanations struct {
}

type NonCharacterTTS struct {
	Tokens map[string]string `json:"tokens"`
	Slow   string            `json:"slow"`
	Normal string            `json:"normal"`
}

type ChallengeToken struct {
	HintTable *HintTable `json:"hintTable,omitempty"`
	TTS       *string    `json:"tts,omitempty"`
	Value     string     `json:"value"`
}

type SessionMetadata struct {
	LanguageString         string                                    `json:"language_string"`
	HintsURL               string                                    `json:"hints_url"`
	Beginner               bool                                      `json:"beginner"`
	MixtureModels          map[string]MixtureModel                   `json:"mixture_models"`
	FirstLesson            bool                                      `json:"first_lesson"`
	PassStrength           float64                                   `json:"pass_strength"`
	LevelSessionNumber     int64                                     `json:"level_session_number"`
	LevelIndex             int64                                     `json:"level_index"`
	SkillID                string                                    `json:"skill_id"`
	MinStrengthIncrement   float64                                   `json:"min_strength_increment"`
	ID                     string                                    `json:"id"`
	FirstUserLesson        bool                                      `json:"first_user_lesson"`
	MinStrengthDecrement   float64                                   `json:"min_strength_decrement"`
	SkillIndex             int64                                     `json:"skill_index"`
	LessonNumber           int64                                     `json:"lesson_number"`
	SkillTreeID            string                                    `json:"skill_tree_id"`
	Type                   string                                    `json:"type"`
	LevelSessionIndex      int64                                     `json:"level_session_index"`
	UILanguage             FromLanguageEnum                          `json:"ui_language"`
	Explanation            string                                    `json:"explanation"`
	SkillTitle             string                                    `json:"skill_title"`
	TeachesLexemeIDS       []string                                  `json:"teaches_lexeme_ids"`
	Language               LearningLanguageEnum                      `json:"language"`
	FromLanguage           FromLanguageEnum                          `json:"from_language"`
	SkillName              string                                    `json:"skill_name"`
	SkillTreeLevel         int64                                     `json:"skill_tree_level"`
	TargetLexemeIDS        []string                                  `json:"target_lexeme_ids"`
	KcStrengthModelVersion int64                                     `json:"kc_strength_model_version"`
	Experiments            []interface{}                             `json:"experiments"`
	TTSEnabled             bool                                      `json:"tts_enabled"`
	SkillColor             string                                    `json:"skill_color"`
	LexemeLevelUpdates     map[string](map[string]LexemeLevelUpdate) `json:"lexeme_level_updates"`
}

type LexemeLevelUpdate struct {
	LevelIndex         int64   `json:"level_index"`
	UntargetedProgress float64 `json:"untargeted_progress"`
	TimesTargeted      int64   `json:"times_targeted"`
	TargetedProgress   float64 `json:"targeted_progress"`
}

type MixtureModel struct {
	Prior          []float64   `json:"prior"`
	LearningCurves [][]float64 `json:"learning_curves"`
}

type ProgressUpdate struct {
	ProgressVersion int64     `json:"progressVersion"`
	Progress        []float64 `json:"progress"`
	LevelIndex      int64     `json:"levelIndex"`
	ProgressKeys    []string  `json:"progressKeys"`
	SkillID         string    `json:"skillId"`
}

type TrackingProperties struct {
	NumChallengesGt                    int64                `json:"num_challenges_gt"`
	LevelSessionIsLast                 bool                 `json:"level_session_is_last"`
	ActivityUUID                       string               `json:"activity_uuid"`
	ReadFromCache                      bool                 `json:"read_from_cache"`
	SumContentLength                   int64                `json:"sum_content_length"`
	NumAdaptiveChallengesGenerated     int64                `json:"num_adaptive_challenges_generated"`
	LevelIndex                         int64                `json:"level_index"`
	SkillID                            string               `json:"skill_id"`
	PercentKcCoverage                  float64              `json:"percent_kc_coverage"`
	DistinctUndirectedSentencesCount   int64                `json:"distinct_undirected_sentences_count"`
	DataVersion                        string               `json:"data_version"`
	NumChallengesGenerated             int64                `json:"num_challenges_generated"`
	LevelSessionIndex                  int64                `json:"level_session_index"`
	NumChallengesGtTap                 int64                `json:"num_challenges_gt_tap"`
	LessonNumber                       int64                `json:"lesson_number"`
	Offline                            bool                 `json:"offline"`
	SkillTreeID                        string               `json:"skill_tree_id"`
	Type                               string               `json:"type"`
	AvgContentLength                   float64              `json:"avg_content_length"`
	LexemesWereReordered               bool                 `json:"lexemes_were_reordered"`
	PercentHighFailureChallenges       float64              `json:"percent_high_failure_challenges"`
	SkillXCoord                        int64                `json:"skill_x_coord"`
	TreeLevel                          int64                `json:"tree_level"`
	GenerationTimestamp                int64                `json:"generation_timestamp"`
	MaxRepeatedSentenceCount           int64                `json:"max_repeated_sentence_count"`
	FromLanguage                       FromLanguageEnum     `json:"from_language"`
	SkillName                          string               `json:"skill_name"`
	ExpectedLength                     int64                `json:"expected_length"`
	MaxRepeatedChallengeCount          int64                `json:"max_repeated_challenge_count"`
	GenerationAppVersion               string               `json:"generation_app_version"`
	NumAdaptiveChallengesGtReverseTap  int64                `json:"num_adaptive_challenges_gt_reverse_tap"`
	LearningLanguage                   LearningLanguageEnum `json:"learning_language"`
	NumChallengesGtListenTap           int64                `json:"num_challenges_gt_listen_tap"`
	SentencesCount                     int64                `json:"sentences_count"`
	IsShorterThanExpected              bool                 `json:"is_shorter_than_expected"`
	NumChallengesGtSpeak               int64                `json:"num_challenges_gt_speak"`
	MaxRepeatedUndirectedSentenceCount int64                `json:"max_repeated_undirected_sentence_count"`
	NumHighFailureChallenges           int64                `json:"num_high_failure_challenges"`
	DistinctSentencesCount             int64                `json:"distinct_sentences_count"`
	NumAdaptiveChallengesGt            int64                `json:"num_adaptive_challenges_gt"`
}

type SpecificType string

const (
	ReverseTap            SpecificType = "reverse_tap"
	SpecificTypeListenTap SpecificType = "listen_tap"
	SpecificTypeSpeak     SpecificType = "speak"
	Tap                   SpecificType = "tap"
)

type FromLanguageEnum string

const (
	En   FromLanguageEnum = "en"
	NlNL FromLanguageEnum = "nl-NL"
)

type VertexType string

const (
	Typo VertexType = "typo"
)

type LearningLanguageEnum string

const (
	DN LearningLanguageEnum = "dn"
)

type KcTypeStr string

const (
	Lex KcTypeStr = "lex"
)

type MetadataType string

const (
	MetaSpeak     MetadataType = "speak"
	MetaTranslate MetadataType = "translate"
	MetaListenTap MetadataType = "listen_tap"
)

type TargetLanguageName string

const (
	English TargetLanguageName = "English"
)

type ChallengeType string

const (
	Speak     ChallengeType = "speak"
	Translate ChallengeType = "translate"
	ListenTap ChallengeType = "listenTap"
)
