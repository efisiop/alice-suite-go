-- Add Alice in Wonderland Character Names to Glossary
-- Book ID: alice-in-wonderland

BEGIN TRANSACTION;

-- Main Characters
INSERT OR IGNORE INTO alice_glossary (id, book_id, term, definition, example, chapter_reference)
VALUES 
('character-alice', 'alice-in-wonderland', 'alice', 'The young protagonist of the story, a curious and imaginative girl who falls down a rabbit hole into Wonderland', 'Alice was beginning to get very tired of sitting by her sister on the bank.', 'I'),

('character-white-rabbit', 'alice-in-wonderland', 'white rabbit', 'A nervous, time-obsessed rabbit who leads Alice into Wonderland. He is always worried about being late.', 'There was nothing so very remarkable in that; nor did Alice think it so very much out of the way to hear the Rabbit say to itself, "Oh dear! Oh dear! I shall be late!"', 'I'),

('character-cheshire-cat', 'alice-in-wonderland', 'cheshire cat', 'A mysterious, grinning cat with the ability to appear and disappear at will. Known for its philosophical remarks and distinctive smile.', 'The Cheshire Cat was now gone, but its grin remained, floating in the air.', 'VI'),

('character-mad-hatter', 'alice-in-wonderland', 'mad hatter', 'A whimsical character who hosts a never-ending tea party. Known for his nonsensical riddles and eccentric behavior.', 'The Hatter opened his eyes very wide on hearing this; but all he said was, "Why is a raven like a writing-desk?"', 'VII'),

('character-march-hare', 'alice-in-wonderland', 'march hare', 'A character who attends the Mad Tea Party. In Victorian England, hares were thought to go mad in March, hence the name.', 'The March Hare took the watch and looked at it gloomily: then he dipped it into his cup of tea.', 'VII'),

('character-queen-of-hearts', 'alice-in-wonderland', 'queen of hearts', 'The tyrannical ruler of Wonderland, known for her frequent orders to behead people. She represents irrational authority.', '"Off with her head!" the Queen shouted at the top of her voice.', 'VIII'),

('character-king-of-hearts', 'alice-in-wonderland', 'king of hearts', 'The husband of the Queen of Hearts, who tries to be reasonable and often pardons people the Queen wants to execute.', 'The King looked anxiously at the White Rabbit, who said in a low voice, "Your Majesty must cross-examine THIS witness."', 'XI'),

('character-duchess', 'alice-in-wonderland', 'duchess', 'A character who appears in Chapter VI. She is initially hostile but later becomes friendly, offering Alice moral lessons.', 'The Duchess was sitting on a three-legged stool in the middle, nursing a baby.', 'VI'),

('character-mock-turtle', 'alice-in-wonderland', 'mock turtle', 'A melancholy character with the head, hind hooves, and tail of a cow. He tells Alice about his school days and sings sad songs.', 'The Mock Turtle sighed deeply, and drew the back of one flapper across his eyes.', 'IX'),

('character-gryphon', 'alice-in-wonderland', 'gryphon', 'A mythical creature with the head and wings of an eagle and the body of a lion. He escorts Alice to meet the Mock Turtle.', 'They very soon came upon a Gryphon, lying fast asleep in the sun.', 'IX'),

('character-mouse', 'alice-in-wonderland', 'mouse', 'A character who appears in Chapter II. Alice meets the Mouse in the Pool of Tears, and he tries to dry everyone off by telling a dry history lesson.', 'The Mouse looked at her rather inquisitively, and seemed to her to wink with one of its little eyes.', 'II'),

('character-caterpillar', 'alice-in-wonderland', 'caterpillar', 'A wise but somewhat rude character who sits on a mushroom smoking a hookah. He gives Alice advice about growing and shrinking.', 'The Caterpillar and Alice looked at each other for some time in silence: at last the Caterpillar took the hookah out of its mouth.', 'V'),

('character-dodo', 'alice-in-wonderland', 'dodo', 'A character who appears in Chapter III. He organizes the Caucus Race, where everyone runs in circles and everyone wins.', 'The Dodo had paused as if it thought that somebody ought to speak, but nobody else spoke.', 'III'),

('character-lory', 'alice-in-wonderland', 'lory', 'A character who appears in Chapter III. The Lory is one of the animals Alice meets in the Pool of Tears.', 'The Lory positively refused to tell its age.', 'III'),

('character-eaglet', 'alice-in-wonderland', 'eaglet', 'A character who appears in Chapter III. The Eaglet is one of the animals Alice meets in the Pool of Tears.', 'The Eaglet bent down its head to hide a smile.', 'III'),

('character-duck', 'alice-in-wonderland', 'duck', 'A character who appears in Chapter III. The Duck is one of the animals Alice meets in the Pool of Tears.', 'The Duck said, "Found it!"', 'III'),

('character-knave-of-hearts', 'alice-in-wonderland', 'knave of hearts', 'A character accused of stealing the Queen''s tarts. He is put on trial in Chapter XI.', 'The Knave of Hearts, he made those tarts, and took them quite away!', 'XI'),

('character-bill-the-lizard', 'alice-in-wonderland', 'bill the lizard', 'A character who appears in Chapter IV. Bill is a lizard who works for the White Rabbit and is sent down the chimney.', 'Bill the Lizard was a footman: and the two footmen were both shaped like the three gardeners.', 'IV');

COMMIT;
