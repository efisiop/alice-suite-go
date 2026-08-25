// Reader UI localization only. Book passages, reader content, dictionary
// definitions, and AI answers remain in their original language.
(function() {
    const originals = new WeakMap();
    const originalAttributePrefix = 'data-alice-i18n-original-';

    const text = {
        'Alice in Wonderland': 'Alice nel Paese delle Meraviglie',
        'Alice Suite Reader': 'Lettore Alice Suite',
        'Home': 'Home',
        'Login': 'Accedi',
        'Log in': 'Accedi',
        'Register': 'Registrati',
        'Create Account': 'Crea account',
        'Create account': 'Crea account',
        'Email': 'Email',
        'Password': 'Password',
        'First Name': 'Nome',
        'Last Name': 'Cognome',
        'Help Language': 'Lingua di aiuto',
        'English': 'Inglese',
        'Italian': 'Italiano',
        'Danish': 'Danese',
        'Spanish': 'Spagnolo',
        'French': 'Francese',
        'German': 'Tedesco',
        'Portuguese': 'Portoghese',
        'Already have an account? Login': 'Hai gia un account? Accedi',
        "Don't have an account? Register": 'Non hai un account? Registrati',
        'JavaScript is disabled.': 'JavaScript e disattivato.',
        'Please enable JavaScript to use the login form.': 'Attiva JavaScript per usare il modulo di accesso.',
        'Logging in...': 'Accesso in corso...',
        'Login failed. Please try again.': 'Accesso non riuscito. Riprova.',
        'Invalid email or password': 'Email o password non validi',
        'Registration failed. Please try again.': 'Registrazione non riuscita. Riprova.',
        'Reader Dashboard': 'Dashboard lettore',
        'Reading': 'Lettura',
        'My Page': 'La mia pagina',
        'Logout': 'Esci',
        'Start Reading': 'Inizia a leggere',
        'Pick up where you left off.': 'Riprendi da dove avevi lasciato.',
        'Open Book': 'Apri libro',
        'Reading Progress': 'Progresso di lettura',
        'See your progress and stats.': 'Vedi i tuoi progressi e le statistiche.',
        'View Stats': 'Vedi statistiche',
        'View your help requests and activity.': 'Vedi richieste di aiuto e attivita.',
        'Go to My Page': 'Vai alla mia pagina',
        'Recent Activity': 'Attivita recente',
        'Loading activity...': 'Caricamento attivita...',
        'No recent activity yet - start reading to see it here.': 'Nessuna attivita recente: inizia a leggere per vederla qui.',
        'No recent activity yet — start reading to see it here.': 'Nessuna attivita recente: inizia a leggere per vederla qui.',
        'My Help Requests': 'Le mie richieste di aiuto',
        'My Messages': 'I miei messaggi',
        'My Progress': 'I miei progressi',
        'Account Settings': 'Impostazioni account',
        'Save Settings': 'Salva impostazioni',
        'Saving...': 'Salvataggio...',
        'Saved': 'Salvato',
        'Could not load settings': 'Impossibile caricare le impostazioni',
        'Save failed': 'Salvataggio non riuscito',
        'Loading...': 'Caricamento...',
        'Error loading': 'Errore di caricamento',
        'Request Help': 'Richiedi aiuto',
        'What do you need help with?': 'Di cosa hai bisogno?',
        'Submit Request': 'Invia richiesta',
        'Submitting...': 'Invio in corso...',
        'No help requests yet': 'Nessuna richiesta di aiuto',
        'Pending': 'In attesa',
        'Assigned': 'Assegnata',
        'Resolved': 'Risolta',
        'Response:': 'Risposta:',
        'Consultant': 'Consulente',
        'No messages yet': 'Nessun messaggio',
        'Consultant responses will appear here': 'Le risposte del consulente appariranno qui',
        'Days Active': 'Giorni attivi',
        'Word Lookups': 'Parole cercate',
        'Total Interactions': 'Interazioni totali',
        'Chapters Progress': 'Progresso capitoli',
        'Complete': 'Completato',
        'Start reading to track progress': 'Inizia a leggere per tracciare i progressi',
        'Reading Statistics': 'Statistiche di lettura',
        'Pages Read': 'Pagine lette',
        'Reading Time': 'Tempo di lettura',
        'Words Looked Up': 'Parole cercate',
        'Loading progress data...': 'Caricamento dati progresso...',
        'Verify Book Code': 'Verifica codice libro',
        'Enter your book verification code to get started reading.': 'Inserisci il codice di verifica del libro per iniziare a leggere.',
        'Verification Code': 'Codice di verifica',
        'Verify': 'Verifica',
        'Book verified successfully!': 'Libro verificato correttamente!',
        'Invalid verification code': 'Codice di verifica non valido',
        'Verification failed. Please try again.': 'Verifica non riuscita. Riprova.',
        'Sections': 'Sezioni',
        "Select the section you're reading": 'Seleziona la sezione che stai leggendo',
        'Loading sections...': 'Caricamento sezioni...',
        'No sections available.': 'Nessuna sezione disponibile.',
        'No sections available. Check console for details.': 'Nessuna sezione disponibile.',
        'No snippets available.': 'Nessun estratto disponibile.',
        'Page': 'Pagina',
        'Section': 'Sezione',
        '- Section': '- Sezione',
        'Jump to table of contents': 'Vai all indice',
        'Services': 'Servizi',
        'Help Center': 'Centro assistenza',
        'Navigation tools': 'Strumenti di navigazione',
        'AI Help': 'Aiuto AI',
        'Open AI Help': 'Apri aiuto AI',
        'Human Consultant': 'Consulente umano',
        'Dictionary': 'Dizionario',
        'Test your knowledge': 'Verifica la comprensione',
        'Info Center': 'Centro informazioni',
        'Ah Ah Moments': 'Momenti ah ah',
        'Lectures in town': 'Conferenze in citta',
        'Publisher connection': 'Contatto editore',
        'Go to Page': 'Vai a pagina',
        'Go': 'Vai',
        'Point camera at page text to find your location': 'Punta la fotocamera sul testo della pagina per trovare la posizione',
        'Scan to Locate Your Reading Position': 'Scansiona per trovare la posizione di lettura',
        "Point your camera at the text on the page you're reading, or upload a screenshot/image of text from": 'Punta la fotocamera sul testo della pagina che stai leggendo oppure carica uno screenshot o un immagine di testo da',
        'Upload image of text': 'Carica immagine del testo',
        'Scan this image': 'Scansiona questa immagine',
        'Align 1-2 lines of text inside the green frame, then tap Capture & Scan': 'Allinea 1-2 righe di testo nella cornice verde, poi tocca Acquisisci e scansiona',
        'Align 1–2 lines of text inside the green frame, then tap Capture & Scan': 'Allinea 1-2 righe di testo nella cornice verde, poi tocca Acquisisci e scansiona',
        "What we're reading (black & white):": 'Cio che stiamo leggendo (bianco e nero):',
        'Processing...': 'Elaborazione...',
        'Processing image...': 'Elaborazione immagine...',
        'Location Found!': 'Posizione trovata!',
        'Error': 'Errore',
        'Start Camera': 'Avvia fotocamera',
        'Capture & Scan': 'Acquisisci e scansiona',
        'Try Again': 'Riprova',
        'Hide sections': 'Nascondi sezioni',
        'Hide services': 'Nascondi servizi',
        'Show sections': 'Mostra sezioni',
        'Show services': 'Mostra servizi',
        'Look up word': 'Cerca parola',
        'Look Up': 'Cerca',
        'Enter word': 'Inserisci parola',
        'Loading definition...': 'Caricamento definizione...',
        'Word not found in dictionary.': 'Parola non trovata nel dizionario.',
        'No word provided': 'Nessuna parola inserita',
        'Error looking up word. Please try again.': 'Errore nella ricerca della parola. Riprova.',
        'Derivation': 'Etimologia',
        'Examples': 'Esempi',
        'Example': 'Esempio',
        'Picture': 'Immagine',
        'Finding picture...': 'Ricerca immagine...',
        'Generating picture...': 'Generazione immagine...',
        'Derivation not available for this word.': 'Etimologia non disponibile per questa parola.',
        'Loading derivation...': 'Caricamento etimologia...',
        'No simple example available for this word.': 'Nessun esempio semplice disponibile per questa parola.',
        'Generated everyday examples': 'Esempi d uso quotidiano generati',
        'Source: Alice glossary': 'Fonte: glossario Alice',
        'Source: dictionary (saved)': 'Fonte: dizionario (salvato)',
        'Source: dictionary': 'Fonte: dizionario',
        'Source: Wikipedia': 'Fonte: Wikipedia',
        'Source: generated image': 'Fonte: immagine generata',
        'AI Assistant': 'Assistente AI',
        'Online': 'Online',
        'Linked to': 'Collegato a',
        'Current section': 'Sezione corrente',
        'Ask about this': 'Chiedi su questo',
        'The assistant will use the section you are reading.': 'L assistente usera la sezione che stai leggendo.',
        'This section': 'Questa sezione',
        'This page': 'Questa pagina',
        'Selected text': 'Testo selezionato',
        'Start with the linked text': 'Inizia dal testo collegato',
        'Pick a useful question, or write your own below.': 'Scegli una domanda utile oppure scrivine una sotto.',
        'Explain what is happening': 'Spiega cosa sta succedendo',
        'Make this easier to read': 'Rendilo piu facile da leggere',
        'What should I notice?': 'Cosa dovrei notare?',
        'Find the misunderstood word': 'Trova la parola non capita',
        'Visual': 'Visuale',
        'Clear': 'Cancella',
        'Text selected': 'Testo selezionato',
        'Add selected text': 'Aggiungi testo selezionato',
        'Tap to open chat': 'Tocca per aprire la chat',
        'Thinking...': 'Sto pensando...',
        'Thinking…': 'Sto pensando...',
        'Sending...': 'Invio in corso...',
        'Send': 'Invia',
        'Loading conversation...': 'Caricamento conversazione...',
        'Error loading conversation': 'Errore nel caricamento della conversazione',
        'Human Consultant Help': 'Aiuto del consulente umano',
        'Share clear details so your consultant can help you faster.': 'Condividi dettagli chiari cosi il consulente puo aiutarti piu rapidamente.',
        'Good request includes:': 'Una buona richiesta include:',
        'What you are trying to understand': 'Cosa stai cercando di capire',
        'Where you are in the book (page and section)': 'Dove sei nel libro (pagina e sezione)',
        'What is confusing right now': 'Cosa ti confonde adesso',
        'Try to write at least 15-20 words.': 'Prova a scrivere almeno 15-20 parole.',
        'What have you already tried? (optional)': 'Cosa hai gia provato? (facoltativo)',
        'Priority': 'Priorita',
        'Normal': 'Normale',
        'High': 'Alta',
        'Urgent': 'Urgente',
        'Cancel': 'Annulla',
        'Send to consultant': 'Invia al consulente',
        'Sending to consultant...': 'Invio al consulente...',
        'Sent successfully. Your consultant will review it.': 'Inviato correttamente. Il tuo consulente lo esaminera.',
        'Could not send request. Please try again.': 'Impossibile inviare la richiesta. Riprova.',
        'Please add a bit more detail (at least 15 characters).': 'Aggiungi qualche dettaglio in piu (almeno 15 caratteri).',
        'Share your successes and realizations with other readers. Others can see moments you mark as shared.': 'Condividi successi e scoperte con altri lettori. Gli altri possono vedere i momenti che scegli di condividere.',
        'Add your moment': 'Aggiungi il tuo momento',
        'What was your ah ah moment?': 'Qual e stato il tuo momento ah ah?',
        'Page (optional)': 'Pagina (facoltativo)',
        'Section (optional)': 'Sezione (facoltativo)',
        'Share with other readers': 'Condividi con altri lettori',
        'Share moment': 'Condividi momento',
        'Moments from you and other readers': 'Momenti tuoi e di altri lettori',
        'No moments yet. Be the first to share one!': 'Ancora nessun momento. Sii il primo a condividerne uno!',
        'Could not load moments. Try again.': 'Impossibile caricare i momenti. Riprova.',
        'Quiz': 'Quiz',
        'Choose what you want to test': 'Scegli cosa vuoi verificare',
        'Choose what to quiz yourself on:': 'Scegli su cosa verificarti:',
        'This section': 'Questa sezione',
        "the part you're reading right now": 'la parte che stai leggendo ora',
        'the full current page': 'l intera pagina corrente',
        "What I've read so far": 'Cio che ho letto finora',
        'from page 1 up to this page': 'da pagina 1 fino a questa pagina',
        'Generate quiz': 'Genera quiz',
        'Generating quiz...': 'Generazione quiz...',
        'Check answer': 'Controlla risposta',
        'This section only': 'Solo questa sezione',
        'This page only': 'Solo questa pagina',
        "Everything you've read so far": 'Tutto cio che hai letto finora',
        'Start quiz': 'Inizia quiz',
        'Next': 'Avanti',
        'Finish': 'Termina',
        'Quiz complete': 'Quiz completato',
        'Try another quiz': 'Prova un altro quiz',
        'Back to reading': 'Torna alla lettura',
        'Back': 'Indietro',
        'Correct! Great job.': 'Corretto! Ottimo lavoro.',
        'No questions were generated. Try a different scope or try again.': 'Non sono state generate domande. Prova un ambito diverso o riprova.',
        'Something went wrong. Please try again.': 'Qualcosa e andato storto. Riprova.',
        'Loading image...': 'Caricamento immagine...',
        'Extracting text from image...': 'Estrazione testo dall immagine...',
        'Processing text...': 'Elaborazione testo...',
        'Finding location in book...': 'Ricerca posizione nel libro...',
        'Need AI assistance?': 'Hai bisogno di aiuto AI?',
        'Selection mode active - choose text in the book': 'Modalita selezione attiva: scegli il testo nel libro',
        'Selection Mode Active - Tap text to select': 'Modalita selezione attiva: tocca il testo da selezionare',
        'Select a passage in the book': 'Seleziona un passaggio nel libro',
        'No readable text is linked yet. Choose a section or select text from the book.': 'Non e ancora collegato alcun testo leggibile. Scegli una sezione o seleziona il testo dal libro.',
        'Select a passage in the book. The selected words will stay highlighted here.': 'Seleziona un passaggio nel libro. Le parole selezionate resteranno evidenziate qui.',
        'All rights reserved.': 'Tutti i diritti riservati.'
    };

    const attrs = {
        'Ask about the linked text...': 'Chiedi del testo collegato...',
        'Type a message to the AI assistant': 'Scrivi un messaggio all assistente AI',
        'Describe what you need help with...': 'Descrivi di cosa hai bisogno...',
        'Enter code': 'Inserisci codice',
        'Page number': 'Numero pagina',
        'Section number': 'Numero sezione',
        'Page': 'Pagina',
        'Section': 'Sezione',
        'Close': 'Chiudi',
        'Close picture': 'Chiudi immagine',
        'Close derivation': 'Chiudi etimologia',
        'Close example': 'Chiudi esempio',
        'Send message': 'Invia messaggio',
        'Hide sections panel': 'Nascondi pannello sezioni',
        'Hide services panel': 'Nascondi pannello servizi',
        'Previous page': 'Pagina precedente',
        'Next page': 'Pagina successiva',
        'Page navigation': 'Navigazione pagine',
        'Reader layout controls': 'Controlli layout lettore',
        'AI suggestion for this passage': 'Suggerimento AI per questo passaggio',
        'AI suggestion': 'Suggerimento AI',
        'Dismiss': 'Ignora',
        'Resize': 'Ridimensiona',
        'Make smaller': 'Riduci',
        'Make larger': 'Ingrandisci',
        'Minimize to bottom': 'Riduci in basso',
        'Drag to resize': 'Trascina per ridimensionare',
        'Need AI assistance': 'Hai bisogno di aiuto AI',
        'Get AI help': 'Ottieni aiuto AI',
        'AI reading context': 'Contesto di lettura AI',
        'AI starter questions': 'Domande iniziali AI',
        'Get visual example': 'Ottieni esempio visivo',
        'Clear conversation': 'Cancella conversazione',
        'Open picture': 'Apri immagine',
        'Click to view full size': 'Fai clic sull immagine per vederla a grandezza naturale',
        'e.g. I finally understood why Alice followed the rabbit...': 'ad es. Ho finalmente capito perche Alice ha seguito il coniglio...',
        "Example: I don't understand why Alice follows the rabbit here, and what this means for the story.": 'Esempio: non capisco perche Alice segue il coniglio qui e cosa significhi per la storia.',
        "Example: I checked the dictionary and AI Help, but I still don't understand the sentence meaning.": 'Esempio: ho controllato il dizionario e l aiuto AI, ma non capisco ancora il significato della frase.',
        'Select text in the book first...': 'Prima seleziona il testo nel libro...',
        'Type a message...': 'Scrivi un messaggio...'
    };

    const patterns = [
        [/^Page (\d+) - Section (\d+)$/, 'Pagina $1 - Sezione $2'],
        [/^Question (\d+) of (\d+)$/, 'Domanda $1 di $2'],
        [/^You got (\d+) out of (\d+) correct\.$/, 'Hai risposto correttamente a $1 domande su $2.'],
        [/^Incorrect\. The correct answer is: (.+)$/, 'Non corretto. La risposta corretta e: $1'],
        [/^Selected: "(.+)"$/, 'Selezionato: "$1"'],
        [/^See (\d+) more example(?:s)?$/, 'Mostra altri $1 esempi'],
        [/^Hide (\d+) more example(?:s)?$/, 'Nascondi altri $1 esempi'],
        [/^Generated everyday examples \((\d+)\/(\d+)\)$/, 'Esempi d uso quotidiano generati ($1/$2)'],
        [/^Quiz on this section \(Section (\d+)\)$/, 'Quiz su questa sezione (Sezione $1)'],
        [/^Quiz on this page \(Page (\d+)\)$/, 'Quiz su questa pagina (Pagina $1)'],
        [/^Quiz on what you've read so far \(Pages 1.(\d+)\)$/, 'Quiz su cio che hai letto finora (Pagine 1-$1)'],
        [/^Picture for (.+)$/, 'Immagine per $1'],
        [/^Generated picture for (.+)$/, 'Immagine generata per $1'],
        [/^Picture not available: (.+)$/, 'Immagine non disponibile: $1'],
        [/^API Error: (.+)\. Please check server logs\.$/, 'Errore API: $1. Controlla i registri del server.'],
        [/^Error looking up word: (.+)\. Please try again\.$/, 'Errore nella ricerca della parola: $1. Riprova.']
    ];

    function normalize(value) {
        return (value || '').replace(/\s+/g, ' ').trim();
    }

    function shouldSkip(node) {
        const element = node.nodeType === Node.ELEMENT_NODE ? node : node.parentElement;
        return !element || !!element.closest('script, style, template, noscript, .page-content, .book-page-content, .reader-book-text, .section-content, .section-snippet, .chat-message-content, .dictionary-popup-definition, .quiz-question-text, .quiz-option');
    }

    function currentLanguage() {
        return (sessionStorage.getItem('alice_reader_language') || 'en').toLowerCase();
    }

    function translateText(value) {
        const normalized = normalize(value);
        if (text[normalized]) return text[normalized];

        for (const [pattern, replacement] of patterns) {
            if (pattern.test(normalized)) return normalized.replace(pattern, replacement);
        }
        return value;
    }

    function translate(value) {
        return currentLanguage() === 'it' ? translateText(value) : value;
    }

    function translateFormat(key, values) {
        let translated = translate(key);
        Object.keys(values || {}).forEach(name => {
            translated = translated.replace(new RegExp('\\{' + name + '\\}', 'g'), values[name]);
        });
        return translated;
    }

    function translateTextNode(node) {
        if (shouldSkip(node)) return;
        const current = node.nodeValue;
        const storedOriginal = originals.get(node);
        const storedTranslation = storedOriginal && translateText(storedOriginal);
        if (!storedOriginal || (current !== storedOriginal && current !== storedTranslation)) {
            originals.set(node, current);
        }
        const original = originals.get(node);
        const translated = translateText(original);
        const nextValue = currentLanguage() === 'it' ? original.replace(normalize(original), translated) : original;
        if (node.nodeValue !== nextValue) node.nodeValue = nextValue;
    }

    function translateTextNodes(root) {
        if (root.nodeType === Node.TEXT_NODE) {
            translateTextNode(root);
            return;
        }
        const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
        while (walker.nextNode()) translateTextNode(walker.currentNode);
    }

    function translateElementAttributes(element) {
        if (element.nodeType !== Node.ELEMENT_NODE || shouldSkip(element)) return;
        const names = ['placeholder', 'aria-label', 'title', 'alt'];
        names.forEach(name => {
            const value = element.getAttribute(name);
            if (!value) return;
            const storageName = originalAttributePrefix + name;
            const storedOriginal = element.getAttribute(storageName);
            const storedTranslation = storedOriginal && attrs[normalize(storedOriginal)];
            const original = !storedOriginal || (value !== storedOriginal && value !== storedTranslation) ? value : storedOriginal;
            if (storedOriginal !== original) element.setAttribute(storageName, original);
            const translated = attrs[normalize(original)];
            const nextValue = currentLanguage() === 'it' && translated ? translated : original;
            if (value !== nextValue) element.setAttribute(name, nextValue);
        });
    }

    function translateAttributes(root) {
        if (root.nodeType === Node.ELEMENT_NODE) translateElementAttributes(root);
        root.querySelectorAll && root.querySelectorAll('*').forEach(translateElementAttributes);
    }

    function translateTree(root) {
        translateTextNodes(root);
        translateAttributes(root);
    }

    let observer = null;
    let appliedLanguage = null;

    function observeMutations() {
        if (observer && document.body) observer.observe(document.body, {
            attributes: true,
            attributeFilter: ['placeholder', 'aria-label', 'title', 'alt'],
            characterData: true,
            childList: true,
            subtree: true
        });
    }

    function apply(languageCode) {
        const language = (languageCode || 'en').toLowerCase();
        sessionStorage.setItem('alice_reader_language', language);
        document.documentElement.lang = language === 'it' ? 'it' : 'en';
        if (!document.body) return;
        if (appliedLanguage === language) return;
        if (observer) observer.disconnect();
        translateTree(document.body);
        appliedLanguage = language;
        observeMutations();
    }

    function loadPreference() {
        const cachedLanguage = sessionStorage.getItem('alice_reader_language');
        if (cachedLanguage) apply(cachedLanguage);

        const token = window.getAuthToken && window.getAuthToken();
        if (!token || !window.location.pathname.startsWith('/reader')) return;
        fetch('/api/reader/preferences', { headers: { Authorization: 'Bearer ' + token } })
            .then(response => response.ok ? response.json() : null)
            .then(preferences => preferences && preferences.preferred_language_code && apply(preferences.preferred_language_code))
            .catch(error => console.warn('[i18n] Could not load reader language preference:', error));
    }

    function watchMutations() {
        if (observer || !document.body) return;
        observer = new MutationObserver(mutations => {
            if (currentLanguage() !== 'it') return;
            observer.disconnect();
            mutations.forEach(mutation => {
                if (mutation.type === 'characterData') {
                    translateTextNode(mutation.target);
                } else if (mutation.type === 'attributes') {
                    translateElementAttributes(mutation.target);
                } else {
                    mutation.addedNodes.forEach(node => translateTree(node));
                }
            });
            observeMutations();
        });
        observeMutations();
    }

    window.aliceReaderTranslations = { it: { text: text, attrs: attrs } };
    window.aliceReaderT = translate;
    window.aliceReaderTFormat = translateFormat;
    window.setAliceReaderLanguagePreference = apply;
    window.loadAliceReaderLanguagePreference = loadPreference;
    window.watchAliceReaderLanguageMutations = watchMutations;
})();
