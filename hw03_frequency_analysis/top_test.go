package hw03frequencyanalysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("additional task:ignore char case", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    []string
			wordHandler func(string) string
		}{
			{
				name:  "English words with strings.ToLower",
				input: "Dog dog dog,",
				expected: []string{
					"dog", "dog,",
				},
				wordHandler: strings.ToLower,
			},
			{
				name:  "Russian words with strings.ToLower",
				input: "Собака собака собака,",
				expected: []string{
					"собака", "собака,",
				},
				wordHandler: strings.ToLower,
			},
			{
				name:  "English words  with helper",
				input: "Dog dog dog,",
				expected: []string{
					"dog", "dog,",
				},
				wordHandler: IgnoreCharCase,
			},
			{
				name:  "Russian words  with helper",
				input: "Собака собака собака,",
				expected: []string{
					"собака", "собака,",
				},
				wordHandler: IgnoreCharCase,
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.expected, Top10(tc.input, tc.wordHandler))
			})
		}
	})

	t.Run("additional task:ignore punctuation characters", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    []string
			wordHandler func(string) string
		}{
			{
				name:  "English words  with helper",
				input: "Dog cat cat,",
				expected: []string{
					"cat", "Dog",
				},
				wordHandler: IgnorePunctuationCharacters,
			},
			{
				name:  "Russian words  with helper",
				input: "Собака кошка кошка,",
				expected: []string{
					"кошка", "Собака",
				},
				wordHandler: IgnorePunctuationCharacters,
			},
			{
				name:  "Russian words  with helper",
				input: "Собака кото-пес котопес,",
				expected: []string{
					"Собака", "кото-пес", "котопес",
				},
				wordHandler: IgnorePunctuationCharacters,
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.expected, Top10(tc.input, tc.wordHandler))
			})
		}
	})

	t.Run("additional task:positive test", func(t *testing.T) {
		expected := []string{
			"а",         // 8
			"он",        // 8
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"в",         // 4
			"его",       // 4
			"если",      // 4
			"кристофер", // 4
			"не",        // 4
		}
		require.Equal(t, expected, Top10(text, IgnoreCharCase, IgnorePunctuationCharacters))
	})

	t.Run("positive test", func(t *testing.T) {
		expected := []string{
			"он",        // 8
			"а",         // 6
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"-",         // 4
			"Кристофер", // 4
			"если",      // 4
			"не",        // 4
			"то",        // 4
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("order words", func(t *testing.T) {
		words := WordCountList{
			WordCountPair{"dog", 5},
			WordCountPair{"parrot", 773},
			WordCountPair{"rabbit", 10},
			WordCountPair{"cat", 5},
			WordCountPair{"eagle", 10},
		}

		wordsSorted := WordCountList{
			WordCountPair{"parrot", 773},
			WordCountPair{"eagle", 10},
			WordCountPair{"rabbit", 10},
			WordCountPair{"cat", 5},
			WordCountPair{"dog", 5},
		}

		require.Equal(t, orderByCountThenByWord(words), wordsSorted)
	})
}
