package solver

import "bytes"
import "errors"
import "encoding/json"

//UnmarshalDuolingoSession ...
func UnmarshalDuolingoSession(data []byte) (DuolingoSession, error) {
	var r DuolingoSession
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DuolingoSession) marshal() ([]byte, error) {
	return json.Marshal(r)
}

//DuolingoSession ...
type DuolingoSession struct {
	FromLanguage             string                  `json:"fromLanguage"`
	LevelSessionIndex        int64                   `json:"levelSessionIndex"`
	LearningLanguage         string                  `json:"learningLanguage"`
	Explanations             Explanations            `json:"explanations"`
	TrackingProperties       TrackingProperties      `json:"trackingProperties"`
	LessonIndex              int64                   `json:"lessonIndex"`
	Metadata                 DuolingoSessionMetadata `json:"metadata"`
	ProgressUpdates          []ProgressUpdate        `json:"progressUpdates"`
	Challenges               []Challenge             `json:"challenges"`
	LevelIndex               int64                   `json:"levelIndex"`
	AdaptiveChallenges       []AdaptiveChallenge     `json:"adaptiveChallenges"`
	SessionStartExperiments  []interface{}           `json:"sessionStartExperiments"`
	Type                     string                  `json:"type"`
	ID                       string                  `json:"id"`
	ChallengeTimeTakenCutoff int64                   `json:"challengeTimeTakenCutoff"`
	SkillID                  string                  `json:"skillId"`
}

//AdaptiveChallenge ...
type AdaptiveChallenge struct {
	SentenceDiscussionID         string                       `json:"sentenceDiscussionId"`
	SentenceID                   string                       `json:"sentenceId"`
	Prompt                       string                       `json:"prompt"`
	SourceLanguage               string                       `json:"sourceLanguage"`
	TTS                          string                       `json:"tts"`
	WrongTokens                  []string                     `json:"wrongTokens"`
	CompactTranslations          []string                     `json:"compactTranslations"`
	ProgressUpdates              []interface{}                `json:"progressUpdates"`
	CorrectIndices               []int64                      `json:"correctIndices"`
	Tokens                       []AdaptiveChallengeToken     `json:"tokens"`
	Grader                       Grader                       `json:"grader"`
	CorrectTokens                []string                     `json:"correctTokens"`
	TargetLanguage               string                       `json:"targetLanguage"`
	ID                           string                       `json:"id"`
	NewWords                     []interface{}                `json:"newWords"`
	ChallengeGeneratorIdentifier ChallengeGeneratorIdentifier `json:"challengeGeneratorIdentifier"`
	CorrectSolutions             []string                     `json:"correctSolutions"`
	Choices                      []ChoiceClass                `json:"choices"`
	Type                         string                       `json:"type"`
	Metadata                     AdaptiveChallengeMetadata    `json:"metadata"`
}

//ChallengeGeneratorIdentifier ...
type ChallengeGeneratorIdentifier struct {
	GeneratorID  string `json:"generatorId"`
	SpecificType string `json:"specificType"`
}

//ChoiceClass ...
type ChoiceClass struct {
	TTS  *string `json:"tts,omitempty"`
	Text string  `json:"text"`
}

//Grader ...
type Grader struct {
	WhitespaceDelimited bool       `json:"whitespaceDelimited"`
	Version             int64      `json:"version"`
	Vertices            [][]Vertex `json:"vertices"`
	Language            string     `json:"language"`
}

//Vertex ...
type Vertex struct {
	To      int64   `json:"to"`
	Lenient string  `json:"lenient"`
	Orig    *string `json:"orig,omitempty"`
	Type    *string `json:"type,omitempty"`
	Auto    *bool   `json:"auto,omitempty"`
}

//AdaptiveChallengeMetadata ...
type AdaptiveChallengeMetadata struct {
	TaggedKcIDS         []TaggedKcID  `json:"tagged_kc_ids"`
	Sentence            string        `json:"sentence"`
	SolutionKey         string        `json:"solution_key"`
	SourceLanguage      string        `json:"source_language"`
	ActivityUUID        string        `json:"activity_uuid"`
	WrongTokens         []string      `json:"wrong_tokens"`
	LexemeIDSToUpdate   []string      `json:"lexeme_ids_to_update"`
	UnknownWords        []interface{} `json:"unknown_words"`
	SpecificType        string        `json:"specific_type"`
	UUID                string        `json:"uuid"`
	WrongLexemes        []string      `json:"wrong_lexemes"`
	GeneratorID         string        `json:"generator_id"`
	Text                string        `json:"text"`
	Type                string        `json:"type"`
	NewExplanationIDS   []interface{} `json:"new_explanation_ids"`
	LexemeIDS           []string      `json:"lexeme_ids"`
	TeachesLexemeIDS    []string      `json:"teaches_lexeme_ids"`
	Tokens              []string      `json:"tokens"`
	SentenceID          string        `json:"sentence_id"`
	DiscussionCount     int64         `json:"discussion_count"`
	TargetLanguageName  string        `json:"target_language_name"`
	Translation         string        `json:"translation"`
	FromLanguage        string        `json:"from_language"`
	HasTTS              bool          `json:"has_tts"`
	LearningLanguage    string        `json:"learning_language"`
	HasAccents          bool          `json:"has_accents"`
	Highlight           []interface{} `json:"highlight"`
	KnowledgeComponents []string      `json:"knowledge_components"`
	TargetLanguage      string        `json:"target_language"`
}

//TaggedKcID ...
type TaggedKcID struct {
	KcTypeStr string `json:"kc_type_str"`
	LegacyID  string `json:"legacy_id"`
}

//AdaptiveChallengeToken ...
type AdaptiveChallengeToken struct {
	HintTable *HintTable `json:"hintTable,omitempty"`
	Value     string     `json:"value"`
}

//HintTable ...
type HintTable struct {
	Headers []string `json:"headers"`
	Rows    [][]Row  `json:"rows"`
}

//Row ...
type Row struct {
	Colspan int64   `json:"colspan"`
	Hint    *string `json:"hint,omitempty"`
}

//Challenge ...
type Challenge struct {
	SentenceDiscussionID         string                       `json:"sentenceDiscussionId"`
	SentenceID                   string                       `json:"sentenceId"`
	Prompt                       string                       `json:"prompt"`
	SourceLanguage               *string                      `json:"sourceLanguage,omitempty"`
	TTS                          *string                      `json:"tts,omitempty"`
	WrongTokens                  []string                     `json:"wrongTokens"`
	CompactTranslations          []string                     `json:"compactTranslations"`
	ProgressUpdates              []interface{}                `json:"progressUpdates"`
	CorrectIndices               []int64                      `json:"correctIndices"`
	Tokens                       []ChallengeToken             `json:"tokens"`
	Grader                       *Grader                      `json:"grader,omitempty"`
	CorrectTokens                []string                     `json:"correctTokens"`
	TargetLanguage               *string                      `json:"targetLanguage,omitempty"`
	ID                           string                       `json:"id"`
	NewWords                     []string                     `json:"newWords"`
	ChallengeGeneratorIdentifier ChallengeGeneratorIdentifier `json:"challengeGeneratorIdentifier"`
	CorrectSolutions             []string                     `json:"correctSolutions"`
	Choices                      []ChoiceUnion                `json:"choices"`
	Type                         string                       `json:"type"`
	Metadata                     ChallengeMetadata            `json:"metadata"`
	SolutionTranslation          *string                      `json:"solutionTranslation,omitempty"`
	SoundID                      *string                      `json:"soundId,omitempty"`
	Threshold                    *float64                     `json:"threshold,omitempty"`
	SlowTTS                      *string                      `json:"slowTts,omitempty"`
	PromptPieces                 []string                     `json:"promptPieces"`
	CorrectIndex                 *int64                       `json:"correctIndex,omitempty"`
}

//ChallengeMetadata ...
type ChallengeMetadata struct {
	TaggedKcIDS         []TaggedKcID      `json:"tagged_kc_ids"`
	NonCharacterTTS     *NonCharacterTTS  `json:"non_character_tts,omitempty"`
	Sentence            *string           `json:"sentence,omitempty"`
	SolutionKey         string            `json:"solution_key"`
	SourceLanguage      string            `json:"source_language"`
	ActivityUUID        string            `json:"activity_uuid"`
	WrongTokens         []string          `json:"wrong_tokens"`
	LexemeIDSToUpdate   []string          `json:"lexeme_ids_to_update"`
	UnknownWords        []string          `json:"unknown_words"`
	SpecificType        string            `json:"specific_type"`
	UUID                string            `json:"uuid"`
	WrongLexemes        []string          `json:"wrong_lexemes"`
	NumComments         int64             `json:"num_comments"`
	GeneratorID         string            `json:"generator_id"`
	Text                string            `json:"text"`
	Type                string            `json:"type"`
	NewExplanationIDS   []interface{}     `json:"new_explanation_ids"`
	GenericLexemeMap    Explanations      `json:"generic_lexeme_map"`
	LexemeIDS           []string          `json:"lexeme_ids"`
	TeachesLexemeIDS    []string          `json:"teaches_lexeme_ids"`
	Tokens              []string          `json:"tokens"`
	SentenceID          string            `json:"sentence_id"`
	DiscussionCount     int64             `json:"discussion_count"`
	TargetLanguageName  *string           `json:"target_language_name,omitempty"`
	Highlight           []string          `json:"highlight"`
	Translation         *string           `json:"translation,omitempty"`
	FromLanguage        string            `json:"from_language"`
	HasTTS              *bool             `json:"has_tts,omitempty"`
	LearningLanguage    string            `json:"learning_language"`
	HasAccents          *bool             `json:"has_accents,omitempty"`
	LexemesToUpdate     []string          `json:"lexemes_to_update"`
	KnowledgeComponents []string          `json:"knowledge_components"`
	TargetLanguage      *string           `json:"target_language,omitempty"`
	Threshold           *float64          `json:"threshold,omitempty"`
	WsLanguage          *string           `json:"ws_language,omitempty"`
	StrippedText        *string           `json:"stripped_text,omitempty"`
	WsThreshold         *float64          `json:"ws_threshold,omitempty"`
	SoundID             *string           `json:"sound_id,omitempty"`
	SolutionTranslation *string           `json:"solution_translation,omitempty"`
	LanguageName        *string           `json:"language_name,omitempty"`
	CorrectSolutions    *CorrectSolutions `json:"correct_solutions"`
	Language            *string           `json:"language,omitempty"`
	NumCorrect          *int64            `json:"num_correct,omitempty"`
	Backoff             *bool             `json:"backoff,omitempty"`
	Sentences           []Option          `json:"sentences"`
	Known               *bool             `json:"known,omitempty"`
	CorrectIndices      []int64           `json:"correct_indices"`
	Options             []Option          `json:"options"`
}

//Explanations ...
type Explanations struct {
}

//NonCharacterTTS ...
type NonCharacterTTS struct {
	Tokens map[string]string `json:"tokens"`
	Slow   string            `json:"slow"`
	Normal string            `json:"normal"`
}

//Option ...
type Option struct {
	I        int    `json:"i"`
	Tabindex int64  `json:"tabindex"`
	Number   int64  `json:"number"`
	Correct  bool   `json:"correct"`
	Sentence string `json:"sentence"`
}

//ChallengeToken ...
type ChallengeToken struct {
	HintTable *HintTable `json:"hintTable,omitempty"`
	TTS       *string    `json:"tts,omitempty"`
	Value     string     `json:"value"`
}

//DuolingoSessionMetadata ...
type DuolingoSessionMetadata struct {
	LanguageString         string                                     `json:"language_string"`
	HintsURL               string                                     `json:"hints_url"`
	Beginner               bool                                       `json:"beginner"`
	MixtureModels          map[string]MixtureModel                    `json:"mixture_models"`
	FirstLesson            bool                                       `json:"first_lesson"`
	PassStrength           float64                                    `json:"pass_strength"`
	LevelSessionNumber     int64                                      `json:"level_session_number"`
	LevelIndex             int64                                      `json:"level_index"`
	SkillID                string                                     `json:"skill_id"`
	MinStrengthIncrement   float64                                    `json:"min_strength_increment"`
	ID                     string                                     `json:"id"`
	FirstUserLesson        bool                                       `json:"first_user_lesson"`
	MinStrengthDecrement   float64                                    `json:"min_strength_decrement"`
	SkillIndex             int64                                      `json:"skill_index"`
	LessonNumber           int64                                      `json:"lesson_number"`
	SkillTreeID            string                                     `json:"skill_tree_id"`
	Type                   string                                     `json:"type"`
	LevelSessionIndex      int64                                      `json:"level_session_index"`
	UILanguage             string                                     `json:"ui_language"`
	SkillTitle             string                                     `json:"skill_title"`
	TeachesLexemeIDS       []string                                   `json:"teaches_lexeme_ids"`
	Language               string                                     `json:"language"`
	FromLanguage           string                                     `json:"from_language"`
	SkillName              string                                     `json:"skill_name"`
	SkillTreeLevel         int64                                      `json:"skill_tree_level"`
	Explanation            string                                     `json:"explanation"`
	TargetLexemeIDS        []string                                   `json:"target_lexeme_ids"`
	KcStrengthModelVersion int64                                      `json:"kc_strength_model_version"`
	Experiments            []interface{}                              `json:"experiments"`
	TTSEnabled             bool                                       `json:"tts_enabled"`
	SkillColor             string                                     `json:"skill_color"`
	LexemeLevelUpdates     map[string](map[string]LexemeLevelUpdates) `json:"lexeme_level_updates"`
}

//LexemeLevelUpdates ...
type LexemeLevelUpdates struct {
	LevelIndex         int64   `json:"level_index"`
	UntargetedProgress float64 `json:"untargeted_progress"`
	TimesTargeted      int64   `json:"times_targeted"`
	TargetedProgress   float64 `json:"targeted_progress"`
}

//MixtureModel ...
type MixtureModel struct {
	Prior          []float64   `json:"prior"`
	LearningCurves [][]float64 `json:"learning_curves"`
}

//ProgressUpdate ...
type ProgressUpdate struct {
	ProgressVersion int64     `json:"progressVersion"`
	Progress        []float64 `json:"progress"`
	LevelIndex      int64     `json:"levelIndex"`
	ProgressKeys    []string  `json:"progressKeys"`
	SkillID         string    `json:"skillId"`
}

//TrackingProperties ...
type TrackingProperties struct {
	NumChallengesGt                    int64   `json:"num_challenges_gt"`
	LevelSessionIsLast                 bool    `json:"level_session_is_last"`
	ActivityUUID                       string  `json:"activity_uuid"`
	ReadFromCache                      bool    `json:"read_from_cache"`
	SumContentLength                   int64   `json:"sum_content_length"`
	NumAdaptiveChallengesGenerated     int64   `json:"num_adaptive_challenges_generated"`
	LevelIndex                         int64   `json:"level_index"`
	SkillID                            string  `json:"skill_id"`
	PercentKcCoverage                  float64 `json:"percent_kc_coverage"`
	DistinctUndirectedSentencesCount   int64   `json:"distinct_undirected_sentences_count"`
	DataVersion                        string  `json:"data_version"`
	NumChallengesGenerated             int64   `json:"num_challenges_generated"`
	LevelSessionIndex                  int64   `json:"level_session_index"`
	NumChallengesGtTap                 int64   `json:"num_challenges_gt_tap"`
	LessonNumber                       int64   `json:"lesson_number"`
	Offline                            bool    `json:"offline"`
	SkillTreeID                        string  `json:"skill_tree_id"`
	Type                               string  `json:"type"`
	AvgContentLength                   float64 `json:"avg_content_length"`
	LexemesWereReordered               bool    `json:"lexemes_were_reordered"`
	PercentHighFailureChallenges       float64 `json:"percent_high_failure_challenges"`
	SkillXCoord                        int64   `json:"skill_x_coord"`
	TreeLevel                          int64   `json:"tree_level"`
	GenerationTimestamp                int64   `json:"generation_timestamp"`
	MaxRepeatedSentenceCount           int64   `json:"max_repeated_sentence_count"`
	FromLanguage                       string  `json:"from_language"`
	SkillName                          string  `json:"skill_name"`
	ExpectedLength                     int64   `json:"expected_length"`
	MaxRepeatedChallengeCount          int64   `json:"max_repeated_challenge_count"`
	GenerationAppVersion               string  `json:"generation_app_version"`
	NumAdaptiveChallengesGtReverseTap  int64   `json:"num_adaptive_challenges_gt_reverse_tap"`
	LearningLanguage                   string  `json:"learning_language"`
	NumChallengesGtListenTap           int64   `json:"num_challenges_gt_listen_tap"`
	SentencesCount                     int64   `json:"sentences_count"`
	IsShorterThanExpected              bool    `json:"is_shorter_than_expected"`
	NumChallengesGtSpeak               int64   `json:"num_challenges_gt_speak"`
	NumChallengesGtTargetLearningJudge int64   `json:"num_challenges_gt_target_learning_judge"`
	MaxRepeatedUndirectedSentenceCount int64   `json:"max_repeated_undirected_sentence_count"`
	NumHighFailureChallenges           int64   `json:"num_high_failure_challenges"`
	DistinctSentencesCount             int64   `json:"distinct_sentences_count"`
	NumAdaptiveChallengesGt            int64   `json:"num_adaptive_challenges_gt"`
}

//ChoiceUnion ...
type ChoiceUnion struct {
	ChoiceClass *ChoiceClass
	String      *string
}

func (x *ChoiceUnion) UnmarshalJSON(data []byte) error {
	x.ChoiceClass = nil
	var c ChoiceClass
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, false, nil, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.ChoiceClass = &c
	}
	return nil
}

func (x *ChoiceUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, false, nil, x.ChoiceClass != nil, x.ChoiceClass, false, nil, false, nil, false)
}

//CorrectSolutions ...
type CorrectSolutions struct {
	IntegerArray []int64
	String       *string
}

func (x *CorrectSolutions) UnmarshalJSON(data []byte) error {
	x.IntegerArray = nil
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.IntegerArray, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *CorrectSolutions) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.IntegerArray != nil, x.IntegerArray, false, nil, false, nil, false, nil, false)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
