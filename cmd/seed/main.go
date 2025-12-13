package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "data/alice-suite.db"

func main() {
	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fmt.Println("🌱 Seeding chapters and data...")

	// Insert Chapter 1 if not exists
	_, err = db.Exec(`
		INSERT OR IGNORE INTO chapters (id, book_id, title, number)
		VALUES ('chapter-1', 'alice-in-wonderland', 'Chapter 1: Down the Rabbit-Hole', 1)
	`)
	if err != nil {
		log.Printf("Warning inserting chapter 1: %v", err)
	}

	// Insert Chapter 2 if not exists
	_, err = db.Exec(`
		INSERT OR IGNORE INTO chapters (id, book_id, title, number)
		VALUES ('chapter-2', 'alice-in-wonderland', 'Chapter 2: The Pool of Tears', 2)
	`)
	if err != nil {
		log.Printf("Warning inserting chapter 2: %v", err)
	}

	// Insert Chapter 3 if not exists
	_, err = db.Exec(`
		INSERT OR IGNORE INTO chapters (id, book_id, title, number)
		VALUES ('chapter-3', 'alice-in-wonderland', 'Chapter 3: A Caucus-Race and a Long Tale', 3)
	`)
	if err != nil {
		log.Printf("Warning inserting chapter 3: %v", err)
	}

	// Chapter 1 Sections
	chapter1Sections := []struct {
		id        string
		title     string
		content   string
		startPage int
		endPage   int
		number    int
	}{
		{
			id:        "chapter-1-section-1",
			title:     "Beginning",
			content:   "Alice was beginning to get very tired of sitting by her sister on the bank, and of having nothing to do: once or twice she had peeped into the book her sister was reading, but it had no pictures or conversations in it, 'and what is the use of a book,' thought Alice 'without pictures or conversations?' So she was considering in her own mind (as well as she could, for the hot day made her feel very sleepy and stupid), whether the pleasure of making a daisy-chain would be worth the trouble of getting up and picking the daisies, when suddenly a White Rabbit with pink eyes ran close by her.",
			startPage: 1,
			endPage:   3,
			number:    1,
		},
		{
			id:        "chapter-1-section-2",
			title:     "The Rabbit",
			content:   "There was nothing so very remarkable in that; nor did Alice think it so very much out of the way to hear the Rabbit say to itself, 'Oh dear! Oh dear! I shall be late!' (when she thought it over afterwards, it occurred to her that she ought to have wondered at this, but at the time it all seemed quite natural); but when the Rabbit actually took a watch out of its waistcoat-pocket, and looked at it, and then hurried on, Alice started to her feet, for it flashed across her mind that she had never before seen a rabbit with either a waistcoat-pocket, or a watch to take out of it, and burning with curiosity, she ran across the field after it, and fortunately was just in time to see it pop down a large rabbit-hole under the hedge.",
			startPage: 4,
			endPage:   6,
			number:    2,
		},
		{
			id:        "chapter-1-section-3",
			title:     "Down the Hole",
			content:   "In another moment down went Alice after it, never once considering how in the world she was to get out again. The rabbit-hole went straight on like a tunnel for some way, and then dipped suddenly down, so suddenly that Alice had not a moment to think about stopping herself before she found herself falling down a very deep well. Either the well was very deep, or she fell very slowly, for she had plenty of time as she went down to look about her and to wonder what was going to happen next. First, she tried to look down and make out what she was coming to, but it was too dark to see anything; then she looked at the sides of the well, and noticed that they were filled with cupboards and book-shelves; here and there she saw maps and pictures hung upon pegs. She took down a jar from one of the shelves as she passed; it was labelled 'ORANGE MARMALADE', but to her great disappointment it was empty: she did not like to drop the jar for fear of killing somebody, so managed to put it into one of the cupboards as she fell past it.",
			startPage: 7,
			endPage:   10,
			number:    3,
		},
		{
			id:        "chapter-1-section-4",
			title:     "The Hall of Doors",
			content:   "Down, down, down. Would the fall never come to an end? 'I wonder how many miles I've fallen by this time?' she said aloud. 'I must be getting somewhere near the centre of the earth. Let me see: that would be four thousand miles down, I think—' (for, you see, Alice had learnt several things of this sort in her lessons in the schoolroom, and though this was not a very good opportunity for showing off her knowledge, as there was no one to listen to her, still it was good practice to say it over) '—yes, that's about the right distance—but then I wonder what Latitude or Longitude I've got to?' (Alice had no idea what Latitude was, or Longitude either, but thought they were nice grand words to say.)",
			startPage: 11,
			endPage:   14,
			number:    4,
		},
		{
			id:        "chapter-1-section-5",
			title:     "The Golden Key",
			content:   "Presently she began again. 'I wonder if I shall fall right through the earth! How funny it'll seem to come out among the people that walk with their heads downward! The Antipathies, I think—' (she was rather glad there was no one listening, this time, as it didn't sound at all the right word) '—but I shall have to ask them what the name of the country is, you know. Please, Ma'am, is this New Zealand or Australia?' (and she tried to curtsey as she spoke—fancy curtseying as you're falling through the air! Do you think you could manage it?) 'And what an ignorant little girl she'll think me for asking! No, it'll never do to ask: perhaps I shall see it written up somewhere.'",
			startPage: 15,
			endPage:   18,
			number:    5,
		},
		{
			id:        "chapter-1-section-6",
			title:     "The Garden Door",
			content:   "Down, down, down. There was nothing else to do, so Alice soon began talking again. 'Dinah'll miss me very much to-night, I should think!' (Dinah was the cat.) 'I hope they'll remember her saucer of milk at tea-time. Dinah my dear! I wish you were down here with me! There are no mice in the air, I'm afraid, but you might catch a bat, and that's very like a mouse, you know. But do cats eat bats, I wonder?' And here Alice began to get rather sleepy, and went on saying to herself, in a dreamy sort of way, 'Do cats eat bats? Do cats eat bats?' and sometimes, 'Do bats eat cats?' for, you see, as she couldn't answer either question, it didn't much matter which way she put it.",
			startPage: 19,
			endPage:   22,
			number:    6,
		},
		{
			id:        "chapter-1-section-7",
			title:     "The Pool of Tears",
			content:   "She felt that she was dozing off, and had just begun to dream that she was walking hand in hand with Dinah, and saying to her very earnestly, 'Now, Dinah, tell me the truth: did you ever eat a bat?' when suddenly, thump! thump! down she came upon a heap of sticks and dry leaves, and the fall was over.",
			startPage: 23,
			endPage:   26,
			number:    7,
		},
	}

	// Insert Chapter 1 sections
	for _, section := range chapter1Sections {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO sections (id, chapter_id, title, content, start_page, end_page, number)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, section.id, "chapter-1", section.title, section.content, section.startPage, section.endPage, section.number)
		if err != nil {
			log.Printf("Warning inserting section %s: %v", section.id, err)
		} else {
			fmt.Printf("✅ Inserted %s\n", section.id)
		}
	}

	// Chapter 2 Sections
	chapter2Sections := []struct {
		id        string
		title     string
		content   string
		startPage int
		endPage   int
		number    int
	}{
		{
			id:        "chapter-2-section-1",
			title:     "Curiouser and Curiouser",
			content:   "'Curiouser and curiouser!' cried Alice (she was so much surprised, that for the moment she quite forgot how to speak good English); 'now I'm opening out like the largest telescope that ever was! Good-bye, feet!' (for when she looked down at her feet, they seemed to be almost out of sight, they were getting so far off). 'Oh, my poor little feet, I wonder who will put on your shoes and stockings for you now, dears? I'm sure I shan't be able! I shall be a great deal too far off to trouble myself about you: you must manage the best way you can; —but I must be kind to them,' thought Alice, 'or perhaps they won't walk the way I want to go! Let me see: I'll give them a new pair of boots every Christmas.'",
			startPage: 27,
			endPage:   30,
			number:    1,
		},
		{
			id:        "chapter-2-section-2",
			title:     "The White Rabbit Again",
			content:   "And she went on planning to herself how she would manage it. 'They must go by the carrier,' she thought; 'and how funny it'll seem, sending presents to one's own feet! And how odd the directions will look!\n\nAlice's Right Foot, Esq.\n  Hearthrug,\n    near the Fender,\n      (with Alice's love).\n\nOh dear, what nonsense I'm talking!'",
			startPage: 31,
			endPage:   34,
			number:    2,
		},
		{
			id:        "chapter-2-section-3",
			title:     "The Hall and the Key",
			content:   "Just then her head struck against the roof of the hall: in fact she was now more than nine feet high, and she at once took up the little golden key and hurried off to the garden door.\n\nPoor Alice! It was as much as she could do, lying down on one side, to look through into the garden with one eye; but to get through was more hopeless than ever: she sat down and began to cry again.",
			startPage: 35,
			endPage:   38,
			number:    3,
		},
		{
			id:        "chapter-2-section-4",
			title:     "The Pool of Tears",
			content:   "'You ought to be ashamed of yourself,' said Alice, 'a great girl like you,' (she might well say this), 'to go on crying in this way! Stop this moment, I tell you!' But she went on all the same, shedding gallons of tears, until there was a large pool all round her, about four inches deep and reaching half down the hall.\n\nAfter a time she heard a little pattering of feet in the distance, and she hastily dried her eyes to see what was coming. It was the White Rabbit returning, splendidly dressed, with a pair of white kid gloves in one hand and a large fan in the other: he came trotting along in a great hurry, muttering to himself as he came, 'Oh! the Duchess, the Duchess! Oh! won't she be savage if I've kept her waiting!'",
			startPage: 39,
			endPage:   42,
			number:    4,
		},
		{
			id:        "chapter-2-section-5",
			title:     "The Fan and Gloves",
			content:   "Alice felt so desperate that she was ready to ask help of any one; so, when the Rabbit came near her, she began, in a low, timid voice, 'If you please, sir—' The Rabbit started violently, dropped the white kid gloves and the fan, and skurried away into the darkness as hard as he could go.\n\nAlice took up the fan and gloves, and, as the hall was very hot, she kept fanning herself all the time she went on talking: 'Dear, dear! How queer everything is to-day! And yesterday things went on just as usual. I wonder if I've been changed in the night? Let me think: was I the same when I got up this morning? I almost think I can remember feeling a little different. But if I'm not the same, the next question is, Who in the world am I? Ah, that's the great puzzle!'",
			startPage: 43,
			endPage:   46,
			number:    5,
		},
		{
			id:        "chapter-2-section-6",
			title:     "The Shrinking",
			content:   "And she began thinking over all the children she knew that were of the same age as herself, to see if she could have been changed for any of them.\n\n'I'm sure I'm not Ada,' she said, 'for her hair goes in such long ringlets, and mine doesn't go in ringlets at all; and I'm sure I can't be Mabel, for I know all sorts of things, and she, oh! she knows such a very little! Besides, she's she, and I'm I, and—oh dear, how puzzling it all is! I'll try if I know all the things I used to know. Let me see: four times five is twelve, and four times six is thirteen, and four times seven is—oh dear! I shall never get to twenty at that rate! However, the Multiplication Table doesn't signify: let's try Geography. London is the capital of Paris, and Paris is the capital of Rome, and Rome—no, that's all wrong, I'm certain! I must have been changed for Mabel! I'll try and say \"How doth the little—\"'",
			startPage: 47,
			endPage:   50,
			number:    6,
		},
		{
			id:        "chapter-2-section-7",
			title:     "The Mouse Appears",
			content:   "She crossed her hands on her lap as if she were saying lessons, and began to repeat it, but her voice sounded hoarse and strange, and the words did not come the same as they used to do:—\n\n'How doth the little crocodile\nImprove his shining tail,\nAnd pour the waters of the Nile\nOn every golden scale!\n\nHow cheerfully he seems to grin,\nHow neatly spread his claws,\nAnd welcome little fishes in\nWith gently smiling jaws!'",
			startPage: 51,
			endPage:   54,
			number:    7,
		},
	}

	// Insert Chapter 2 sections
	for _, section := range chapter2Sections {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO sections (id, chapter_id, title, content, start_page, end_page, number)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, section.id, "chapter-2", section.title, section.content, section.startPage, section.endPage, section.number)
		if err != nil {
			log.Printf("Warning inserting section %s: %v", section.id, err)
		} else {
			fmt.Printf("✅ Inserted %s\n", section.id)
		}
	}

	// Chapter 3 Sections
	chapter3Sections := []struct {
		id        string
		title     string
		content   string
		startPage int
		endPage   int
		number    int
	}{
		{
			id:        "chapter-3-section-1",
			title:     "The Mouse",
			content:   "'O Mouse, do you know the way out of this pool? I am very tired of swimming about here, O Mouse!' (Alice thought this must be the right way of speaking to a mouse: she had never done such a thing before, but she remembered having seen in her brother's Latin Grammar, 'A mouse—of a mouse—to a mouse—a mouse—O mouse!') The Mouse looked at her rather inquisitively, and seemed to her to wink with one of its little eyes, but it said nothing.\n\n'Perhaps it doesn't understand English,' thought Alice; 'I daresay it's a French mouse, come over with William the Conqueror.' (For, with all her knowledge of history, Alice had no very clear notion how long ago anything had happened.) So she began again: 'Où est ma chatte?' which was the first sentence in her French lesson-book. The Mouse gave a sudden leap out of the water, and seemed to quiver all over with fright. 'Oh, I beg your pardon!' cried Alice hastily, afraid that she had hurt the poor animal's feelings. 'I quite forgot you didn't like cats.'",
			startPage: 55,
			endPage:   58,
			number:    1,
		},
		{
			id:        "chapter-3-section-2",
			title:     "The Mouse's Tale",
			content:   "'Not like cats!' cried the Mouse, in a shrill, passionate voice. 'Would you like cats if you were me?'\n\n'Well, perhaps not,' said Alice in a soothing tone: 'don't be angry about it. And yet I wish I could show you our cat Dinah: I think you'd take a fancy to cats if you could only see her. She is such a dear quiet thing,' Alice went on, half to herself, as she swam lazily about in the pool, 'and she sits purring so nicely by the fire, licking her paws and washing her face—and she is such a nice soft thing to nurse—and she's such a capital one for catching mice—oh, I beg your pardon!' cried Alice again, for this time the Mouse was bristling all over, and she felt certain it must be really offended. 'We won't talk about her any more if you'd rather not.'",
			startPage: 59,
			endPage:   62,
			number:    2,
		},
		{
			id:        "chapter-3-section-3",
			title:     "The Caucus-Race",
			content:   "'We indeed!' cried the Mouse, who was trembling down to the end of his tail. 'As if I would talk on such a subject! Our family always hated cats: nasty, low, vulgar things! Don't let me hear the name again!'\n\n'I won't indeed!' said Alice, in a great hurry to change the subject of conversation. 'Are you—are you fond—of—of dogs?' The Mouse did not answer, so Alice went on eagerly: 'There is such a nice little dog near our house I should like to show you! A little bright-eyed terrier, you know, with oh, such long curly brown hair! And it'll fetch things when you throw them, and it'll sit up and beg for its dinner, and all sorts of things—I can't remember half of them—and it belongs to a farmer, you know, and he says it's so useful, it's worth a hundred pounds! He says it kills all the rats and—oh dear!' cried Alice in a sorrowful tone, 'I'm afraid I've offended it again!' For the Mouse was swimming away from her as hard as it could go, and making quite a commotion in the pool as it went.",
			startPage: 63,
			endPage:   66,
			number:    3,
		},
		{
			id:        "chapter-3-section-4",
			title:     "The Dodo's Plan",
			content:   "So she called softly after it, 'Mouse dear! Do come back again, and we won't talk about cats or dogs either, if you don't like them!' When the Mouse heard this, it turned round and swam slowly back to her: its face was quite pale (with passion, Alice thought), and it said in a low trembling voice, 'Let us get to the shore, and then I'll tell you my history, and you'll understand why it is I hate cats and dogs.'\n\nIt was high time to go, for the pool was getting quite crowded with the birds and animals that had fallen into it: there were a Duck and a Dodo, a Lory and an Eaglet, and several other curious creatures. Alice led the way and the whole party swam to the shore.",
			startPage: 67,
			endPage:   70,
			number:    4,
		},
		{
			id:        "chapter-3-section-5",
			title:     "The Race Begins",
			content:   "They were indeed a queer-looking party that assembled on the bank—the birds with draggled feathers, the animals with their fur clinging close to them, and all dripping wet, cross, and uncomfortable.\n\nThe first question of course was, how to get dry again: they had a consultation about this, and after a few minutes it seemed quite natural to Alice to find herself talking familiarly with them, as if she had known them all her life. Indeed, she had quite a long argument with the Lory, who at last turned sulky, and would only say, 'I am older than you, and must know better'; and this Alice would not allow without knowing how old it was, and, as the Lory positively refused to tell its age, there was no more to be said.",
			startPage: 71,
			endPage:   74,
			number:    5,
		},
		{
			id:        "chapter-3-section-6",
			title:     "The Race Results",
			content:   "At last the Mouse, who seemed to be a person of authority among them, called out, 'Sit down, all of you, and listen to me! I'll soon make you dry enough!' They all sat down at once, in a large ring, with the Mouse in the middle. Alice kept her eyes anxiously fixed on it, for she felt sure she would catch a bad cold if she did not get dry very soon.\n\n'Ahem!' said the Mouse with an important air, 'are you all ready? This is the driest thing I know. Silence all round, if you please! \"William the Conqueror, whose cause was favoured by the pope, was soon submitted to by the English, who wanted leaders, and had been of late much accustomed to usurpation and conquest. Edwin and Morcar, the earls of Mercia and Northumbria—\"'\n\n'Ugh!' said the Lory, with a shiver.",
			startPage: 75,
			endPage:   78,
			number:    6,
		},
		{
			id:        "chapter-3-section-7",
			title:     "The Prizes",
			content:   "'I beg your pardon!' said the Mouse, frowning, but very politely: 'Did you speak?'\n\n'Not I!' said the Lory hastily.\n\n'I thought you did,' said the Mouse. '—I proceed. \"Edwin and Morcar, the earls of Mercia and Northumbria, declared for him: and even Stigand, the patriotic archbishop of Canterbury, found it advisable—\"'\n\n'Found what?' said the Duck.\n\n'Found it,' the Mouse replied rather crossly: 'of course you know what \"it\" means.'\n\n'I know what \"it\" means well enough, when I find a thing,' said the Duck: 'it's generally a frog or a worm. The question is, what did the archbishop find?'",
			startPage: 79,
			endPage:   82,
			number:    7,
		},
	}

	// Insert Chapter 3 sections
	for _, section := range chapter3Sections {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO sections (id, chapter_id, title, content, start_page, end_page, number)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, section.id, "chapter-3", section.title, section.content, section.startPage, section.endPage, section.number)
		if err != nil {
			log.Printf("Warning inserting section %s: %v", section.id, err)
		} else {
			fmt.Printf("✅ Inserted %s\n", section.id)
		}
	}

	// Insert verification codes
	codes := []string{"ALICE123", "WONDERLAND", "RABBIT"}
	for _, code := range codes {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO verification_codes (code, book_id, is_used)
			VALUES (?, 'alice-in-wonderland', 0)
		`, code)
		if err != nil {
			log.Printf("Warning inserting code %s: %v", code, err)
		} else {
			fmt.Printf("✅ Inserted verification code: %s\n", code)
		}
	}

	// Insert glossary terms
	glossaryTerms := []struct {
		id              string
		term            string
		definition      string
		chapterRef      string
		example         string
	}{
		{
			id:         "gloss-1",
			term:       "curiouser",
			definition:  "More curious; a playful, non-standard form of the word \"curious\"",
			chapterRef: "chapter-2",
			example:    "Curiouser and curiouser!",
		},
		{
			id:         "gloss-2",
			term:       "waistcoat-pocket",
			definition:  "A small pocket in a waistcoat (vest), a Victorian-era garment",
			chapterRef: "chapter-1",
			example:    "The Rabbit took a watch out of its waistcoat-pocket",
		},
		{
			id:         "gloss-3",
			term:       "caucus-race",
			definition:  "A nonsensical race where everyone runs in circles and everyone wins",
			chapterRef: "chapter-3",
			example:    "A Caucus-Race and a Long Tale",
		},
	}

	for _, term := range glossaryTerms {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO alice_glossary (id, book_id, term, definition, chapter_reference, example)
			VALUES (?, 'alice-in-wonderland', ?, ?, ?, ?)
		`, term.id, term.term, term.definition, term.chapterRef, term.example)
		if err != nil {
			log.Printf("Warning inserting glossary term %s: %v", term.term, err)
		} else {
			fmt.Printf("✅ Inserted glossary term: %s\n", term.term)
		}
	}

	fmt.Println("\n🎉 Seeding completed!")
}

