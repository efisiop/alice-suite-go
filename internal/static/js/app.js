// Alice Suite Reader - Main JavaScript

// Auth helper functions
// Using sessionStorage instead of localStorage to ensure each browser tab/window
// has its own isolated token storage, preventing session mixing when multiple
// readers log in from the same IP address
function getAuthToken() {
    return sessionStorage.getItem('auth_token');
}

function setAuthToken(token) {
    sessionStorage.setItem('auth_token', token);
    // Also sync to cookie for server-side page navigation
    syncTokenToCookie(token);
}

function removeAuthToken() {
    sessionStorage.removeItem('auth_token');
    // Also clear cookie
    document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=Lax';
}

// Sync token from sessionStorage to cookie (for server-side navigation)
// Safari-compatible cookie setting
function syncTokenToCookie(token) {
    if (!token) return;
    const expires = new Date();
    expires.setTime(expires.getTime() + (24 * 60 * 60 * 1000)); // 24 hours
    
    // Safari requires explicit cookie format - use multiple methods for compatibility
    // Method 1: Standard cookie string
    document.cookie = `auth_token=${encodeURIComponent(token)}; expires=${expires.toUTCString()}; path=/; SameSite=Lax`;
    
    // Method 2: Also try without encoding (Safari sometimes prefers this)
    // Note: We encode to be safe, but Safari might need the raw value
    try {
        // Verify cookie was set by reading it back
        const cookies = document.cookie.split(';');
        let found = false;
        for (let cookie of cookies) {
            const [name, value] = cookie.trim().split('=');
            if (name === 'auth_token') {
                found = true;
                break;
            }
        }
        // If not found, try alternative method (Safari-specific)
        if (!found && /Safari/.test(navigator.userAgent) && !/Chrome/.test(navigator.userAgent)) {
            // Safari-specific: try setting without encoding
            document.cookie = `auth_token=${token}; expires=${expires.toUTCString()}; path=/; SameSite=Lax; Secure=false`;
        }
    } catch (e) {
        console.warn('Cookie sync warning:', e);
    }
}

// Ensure cookie is synced from sessionStorage on page load
// Use a flag to prevent multiple simultaneous executions
let cookieSyncInProgress = false;
function ensureCookieSync() {
    // Prevent multiple simultaneous executions
    if (cookieSyncInProgress) {
        return;
    }
    cookieSyncInProgress = true;
    
    try {
        const token = sessionStorage.getItem('auth_token');
        if (token) {
            // Check if cookie exists and matches (Safari-compatible check)
            const cookies = document.cookie.split(';');
            let cookieExists = false;
            let cookieValue = null;
            
            for (let cookie of cookies) {
                const parts = cookie.trim().split('=');
                const name = parts[0];
                const value = parts.slice(1).join('='); // Handle values with = in them
                
                if (name === 'auth_token') {
                    cookieValue = decodeURIComponent(value);
                    // Compare both encoded and decoded versions for Safari compatibility
                    if (cookieValue === token || value === token) {
                        cookieExists = true;
                        break;
                    }
                }
            }
            
            // If cookie doesn't exist or doesn't match, sync it
            if (!cookieExists || cookieValue !== token) {
                syncTokenToCookie(token);
                
                // Safari-specific: verify cookie was set after a short delay
                if (/Safari/.test(navigator.userAgent) && !/Chrome/.test(navigator.userAgent)) {
                    setTimeout(() => {
                        const cookiesAfter = document.cookie.split(';');
                        let verified = false;
                        for (let cookie of cookiesAfter) {
                            const [name, value] = cookie.trim().split('=');
                            if (name === 'auth_token') {
                                const decoded = decodeURIComponent(value);
                                if (decoded === token || value === token) {
                                    verified = true;
                                    break;
                                }
                            }
                        }
                        if (!verified) {
                            console.warn('Safari: Cookie sync may have failed, retrying...');
                            syncTokenToCookie(token);
                        }
                    }, 200);
                }
            }
        }
    } finally {
        // Reset flag after a short delay to allow normal operation
        setTimeout(() => {
            cookieSyncInProgress = false;
        }, 100);
    }
}

function isAuthenticated() {
    return !!getAuthToken();
}

// Global logout function for reader app
// Note: Consultant dashboard defines its own logout function in the head, which takes precedence
window.logout = async function() {
    console.log('[logout] Logout function called');
    
    // Check if we're on consultant dashboard - if so, don't override the consultant logout
    if (window.isConsultantDashboard || window.location.pathname.startsWith('/consultant')) {
        // Consultant dashboard has its own logout function, don't override it
        console.log('[logout] Consultant dashboard detected, using consultant logout');
        // Call the consultant logout if it exists
        if (typeof window.consultantLogout === 'function') {
            window.consultantLogout();
        }
        return;
    }
    
    // Hide user info immediately
    const userInfoNav = document.getElementById('user-info-nav');
    if (userInfoNav) {
        userInfoNav.style.display = 'none';
    }
    
    // Hide name in navbar brand as well
    const userNameBrand = document.getElementById('user-name-brand');
    if (userNameBrand) {
        userNameBrand.style.display = 'none';
        userNameBrand.textContent = '';
    }
    
    // IMPORTANT: Call the logout API BEFORE removing the token
    // This ensures the server records the logout and broadcasts to consultants
    const token = getAuthToken();
    if (token) {
        try {
            console.log('[logout] Calling logout API...');
            const response = await fetch('/auth/v1/logout', {
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            });
            if (response.ok) {
                console.log('[logout] Logout API call successful');
            } else {
                console.warn('[logout] Logout API returned:', response.status);
            }
        } catch (e) {
            console.error('[logout] Error calling logout API:', e);
            // Continue with local logout even if API fails
        }
    }
    
    // Remove auth token
    removeAuthToken();
    
    // Close SSE connection if exists
    if (window.sseConnection) {
        try {
            window.sseConnection.close();
        } catch(e) {
            console.error('[logout] Error closing window.sseConnection:', e);
        }
        window.sseConnection = null;
    }
    if (typeof sseConnection !== 'undefined' && sseConnection) {
        try {
            sseConnection.close();
        } catch(e) {
            console.error('[logout] Error closing sseConnection:', e);
        }
        sseConnection = null;
    }
    
    // Disconnect SSE using the disconnect function if available
    if (typeof disconnectSSE === 'function') {
        try {
            disconnectSSE();
        } catch(e) {
            console.error('[logout] Error calling disconnectSSE:', e);
        }
    }
    
    // Redirect to reader login page
    console.log('[logout] Redirecting to reader login...');
    window.location.replace('/reader/login');
};

// Use event delegation for logout links (works even if links are added dynamically)
document.addEventListener('click', function(e) {
    const logoutLink = e.target.closest('#logout-link-reader');
    if (logoutLink) {
        e.preventDefault();
        console.log('[event delegation] Logout link clicked via delegation');
        if (window.logout) {
            window.logout();
        } else {
            console.error('[event delegation] window.logout not available, using fallback');
            removeAuthToken();
            window.location.href = '/reader/login';
        }
    }
});

// API helper functions
function apiRequest(url, options = {}) {
    const token = getAuthToken();
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };
    
    if (token) {
        headers['Authorization'] = 'Bearer ' + token;
    }
    
    return fetch(url, {
        ...options,
        headers
    });
}

// Reader interface localization. This is intentionally UI-only: it skips the
// book/page text containers so the physical book source stays unchanged.
const aliceReaderTranslations = {
    it: {
        text: {
            'Alice in Wonderland': 'Alice nel Paese delle Meraviglie',
            'Home': 'Home',
            'Login': 'Accedi',
            'Register': 'Registrati',
            'Create Account': 'Crea account',
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
            'Logging in...': 'Accesso in corso...',
            'Login failed. Please try again.': 'Accesso non riuscito. Riprova.',
            'Invalid email or password': 'Email o password non validi',
            'Reading': 'Lettura',
            'My Page': 'La mia pagina',
            'Logout': 'Esci',
            'Reader Dashboard': 'Dashboard lettore',
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
            'Request Help': 'Richiedi aiuto',
            'What do you need help with?': 'Di cosa hai bisogno?',
            'Submit Request': 'Invia richiesta',
            'No help requests yet': 'Nessuna richiesta di aiuto',
            'Pending': 'In attesa',
            'Assigned': 'Assegnata',
            'Resolved': 'Risolta',
            'No messages yet': 'Nessun messaggio',
            'Consultant responses will appear here': 'Le risposte del consulente appariranno qui',
            'Days Active': 'Giorni attivi',
            'Word Lookups': 'Parole cercate',
            'Total Interactions': 'Interazioni totali',
            'Chapters Progress': 'Progresso capitoli',
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
            'Sections': 'Sezioni',
            "Select the section you're reading": 'Seleziona la sezione che stai leggendo',
            'Loading sections...': 'Caricamento sezioni...',
            'Page': 'Pagina',
            ' - Section ': ' - Sezione ',
            'Section': 'Sezione',
            'Jump to table of contents': 'Vai all indice',
            'Services': 'Servizi',
            'AI Help': 'Aiuto AI',
            'Human Consultant': 'Consulente umano',
            'Dictionary': 'Dizionario',
            'Test your knowledge': 'Verifica la comprensione',
            'Info Center': 'Centro informazioni',
            'Ah Ah Moments': 'Momenti ah ah',
            'Lectures in town': 'Conferenze in citta',
            'Publisher connection': 'Contatto editore',
            'Go to Page': 'Vai a pagina',
            'Point camera at page text to find your location': 'Punta la fotocamera sul testo della pagina per trovare la posizione',
            'Hide sections': 'Nascondi sezioni',
            'Hide services': 'Nascondi servizi',
            'Show sections': 'Mostra sezioni',
            'Show services': 'Mostra servizi',
            'Look up word': 'Cerca parola',
            'Look Up': 'Cerca',
            'Enter word': 'Inserisci parola',
            'Loading definition...': 'Caricamento definizione...',
            'Word not found in dictionary.': 'Parola non trovata nel dizionario.',
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
            'Human Consultant Help': 'Aiuto del consulente umano',
            'Share clear details so your consultant can help you faster.': 'Condividi dettagli chiari cosi il consulente puo aiutarti piu rapidamente.',
            'Good request includes:': 'Una buona richiesta include:',
            'What you are trying to understand': 'Cosa stai cercando di capire',
            'Where you are in the book (page and section)': 'Dove sei nel libro (pagina e sezione)',
            'What is confusing right now': 'Cosa ti confonde adesso',
            'Priority': 'Priorita',
            'Normal': 'Normale',
            'High': 'Alta',
            'Urgent': 'Urgente',
            'Cancel': 'Annulla',
            'Send to consultant': 'Invia al consulente',
            'Share with other readers': 'Condividi con altri lettori',
            'Share moment': 'Condividi momento',
            'Moments from you and other readers': 'Momenti tuoi e di altri lettori',
            'No moments yet. Be the first to share one!': 'Ancora nessun momento. Sii il primo a condividerne uno!'
        },
        attrs: {
            'Ask about the linked text...': 'Chiedi del testo collegato...',
            'Type a message to the AI assistant': 'Scrivi un messaggio all assistente AI',
            'Describe what you need help with...': 'Descrivi di cosa hai bisogno...',
            'Enter code': 'Inserisci codice',
            'Page number': 'Numero pagina',
            'Section number': 'Numero sezione',
            "Example: I checked the dictionary and AI Help, but I still don't understand the sentence meaning.": 'Esempio: ho controllato il dizionario e l aiuto AI, ma non capisco ancora il significato della frase.',
            'Close': 'Chiudi',
            'Send message': 'Invia messaggio',
            'Hide sections panel': 'Nascondi pannello sezioni',
            'Hide services panel': 'Nascondi pannello servizi',
            'Previous page': 'Pagina precedente',
            'Next page': 'Pagina successiva',
            'AI suggestion for this passage': 'Suggerimento AI per questo passaggio',
            'Dismiss': 'Ignora',
            'Resize': 'Ridimensiona',
            'Minimize to bottom': 'Riduci in basso',
            'Drag to resize': 'Trascina per ridimensionare'
        }
    }
};

function normalizeI18nText(text) {
    return text.replace(/\s+/g, ' ').trim();
}

function shouldSkipI18nNode(node) {
    const element = node.nodeType === Node.ELEMENT_NODE ? node : node.parentElement;
    if (!element) return true;
    return !!element.closest('script, style, template, noscript, #page-content, #section-snippets, .section-snippet, .page-content, .book-page-content, .reader-book-text, .chat-message-content, .dictionary-popup-definition');
}

function translateTextNodes(root, translations) {
    const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
        acceptNode(node) {
            if (shouldSkipI18nNode(node)) return NodeFilter.FILTER_REJECT;
            const normalized = normalizeI18nText(node.nodeValue || '');
            return translations[normalized] ? NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_SKIP;
        }
    });
    const nodes = [];
    while (walker.nextNode()) {
        nodes.push(walker.currentNode);
    }
    nodes.forEach(node => {
        const original = node.nodeValue || '';
        const normalized = normalizeI18nText(original);
        node.nodeValue = original.replace(normalized, translations[normalized]);
    });
}

function translateAttributes(root, attrTranslations) {
    const attrs = ['placeholder', 'aria-label', 'title'];
    root.querySelectorAll('*').forEach(el => {
        if (shouldSkipI18nNode(el)) return;
        attrs.forEach(attr => {
            const value = el.getAttribute(attr);
            if (!value) return;
            const normalized = normalizeI18nText(value);
            if (attrTranslations[normalized]) {
                el.setAttribute(attr, attrTranslations[normalized]);
            }
        });
    });
}

function applyAliceReaderLanguage(languageCode) {
    const language = (languageCode || 'en').toLowerCase();
    document.documentElement.lang = language === 'it' ? 'it' : 'en';
    sessionStorage.setItem('alice_reader_language', language);

    if (language !== 'it') return;
    const pack = aliceReaderTranslations.it;
    translateTextNodes(document.body, pack.text);
    translateAttributes(document.body, pack.attrs);
}

window.setAliceReaderLanguagePreference = applyAliceReaderLanguage;

function loadAliceReaderLanguagePreference() {
    const cachedLanguage = sessionStorage.getItem('alice_reader_language');
    if (cachedLanguage) {
        applyAliceReaderLanguage(cachedLanguage);
    }

    const token = getAuthToken();
    const path = window.location.pathname;
    if (!token || !path.startsWith('/reader')) return;

    fetch('/api/reader/preferences', {
        headers: {'Authorization': 'Bearer ' + token}
    })
    .then(res => res.ok ? res.json() : null)
    .then(preferences => {
        if (preferences && preferences.preferred_language_code) {
            applyAliceReaderLanguage(preferences.preferred_language_code);
        }
    })
    .catch(err => console.warn('[i18n] Could not load reader language preference:', err));
}

function watchAliceReaderLanguageMutations() {
    let scheduled = false;
    const observer = new MutationObserver(() => {
        const language = sessionStorage.getItem('alice_reader_language');
        if (language !== 'it' || scheduled) return;
        scheduled = true;
        window.setTimeout(() => {
            scheduled = false;
            applyAliceReaderLanguage(language);
        }, 50);
    });
    observer.observe(document.body, { childList: true, subtree: true });
}

// Dictionary lookup
function lookupWord(word, bookId, sectionId) {
    return apiRequest('/rest/v1/rpc/get_definition_with_context', {
        method: 'POST',
        body: JSON.stringify({
            term: word,
            book_id: bookId,
            section_id: sectionId
        })
    }).then(res => res.json());
}

// Show dictionary popup
function showDictionaryPopup(word, definition, x, y) {
    // Remove existing popup
    const existing = document.getElementById('dictionary-popup');
    if (existing) {
        existing.remove();
    }
    
    // Create popup
    const popup = document.createElement('div');
    popup.id = 'dictionary-popup';
    popup.className = 'dictionary-popup';
    popup.style.left = x + 'px';
    popup.style.top = y + 'px';
    popup.innerHTML = `
        <strong>${word}</strong>
        <p>${definition}</p>
        <button class="btn btn-sm btn-secondary" onclick="this.parentElement.remove()">Close</button>
    `;
    
    document.body.appendChild(popup);
    
    // Remove on click outside
    setTimeout(() => {
        document.addEventListener('click', function removePopup(e) {
            if (!popup.contains(e.target)) {
                popup.remove();
                document.removeEventListener('click', removePopup);
            }
        });
    }, 100);
}

// SSE (Server-Sent Events) connection for real-time updates
let sseConnection = null;

function connectSSE() {
    // Don't connect if we're on consultant dashboard
    if (window.isConsultantDashboard || window.location.pathname.indexOf('/consultant') !== -1) {
        return;
    }
    
    const token = getAuthToken();
    if (!token) {
        return;
    }

    // Close existing connection
    if (sseConnection) {
        sseConnection.close();
    }

    // Create new SSE connection
    const eventSource = new EventSource(`/api/realtime/events?token=${encodeURIComponent(token)}`);
    
    eventSource.onmessage = function(event) {
        try {
            if (!event.data || event.data.trim() === '') {
                console.warn('[SSE] Received empty event data');
                return;
            }
            const data = JSON.parse(event.data);
            handleSSEEvent(data);
        } catch (error) {
            console.error('[SSE] Error parsing event data:', error);
            console.error('[SSE] Event data that failed to parse:', event.data);
            // Don't throw - continue processing other events
        }
    };

    eventSource.onerror = function(error) {
        // Only log error if not on consultant dashboard (to avoid noise)
        if (!window.isConsultantDashboard && window.location.pathname.indexOf('/consultant') === -1) {
            console.error('SSE connection error:', error);
            // Reconnect after 5 seconds
            setTimeout(connectSSE, 5000);
        }
    };

    sseConnection = eventSource;
}

function disconnectSSE() {
    if (sseConnection) {
        sseConnection.close();
        sseConnection = null;
    }
}

function handleSSEEvent(event) {
    switch (event.type) {
        case 'help_request':
        case 'help_request_update':
            // Refresh help requests list
            if (typeof refreshHelpRequests === 'function') {
                refreshHelpRequests();
            }
            break;
        case 'activity':
            // Update activity feed
            if (typeof updateActivityFeed === 'function') {
                updateActivityFeed(event.data);
            }
            break;
        case 'online_users':
            // Update online users list
            if (typeof updateOnlineUsers === 'function') {
                updateOnlineUsers(event.data);
            }
            break;
        case 'login':
        case 'logout':
            // Update online users count
            if (typeof updateOnlineUsersCount === 'function') {
                updateOnlineUsersCount();
            }
            break;
    }
}

// Load and display user info in navbar (for reader app)
function loadUserInfoInNavbar() {
    const userInfoNav = document.getElementById('user-info-nav');
    const userNameDisplay = document.getElementById('user-name-display');
    
    // Always hide user info first (in case of logout)
    if (userInfoNav) {
        userInfoNav.style.display = 'none';
    }
    
    const token = getAuthToken();
    if (!token) {
        console.log('[loadUserInfoInNavbar] No auth token found, hiding user info');
        return;
    }

    console.log('[loadUserInfoInNavbar] Elements found:', {
        userInfoNav: !!userInfoNav,
        userNameDisplay: !!userNameDisplay
    });

    console.log('[loadUserInfoInNavbar] Fetching user info from /auth/v1/user');
    fetch('/auth/v1/user', {
        headers: {'Authorization': 'Bearer ' + token}
    })
    .then(res => {
        console.log('[loadUserInfoInNavbar] Response status:', res.status);
        if (!res.ok) {
            throw new Error('Failed to fetch user info: ' + res.status);
        }
        return res.json();
    })
    .then(user => {
        console.log('[loadUserInfoInNavbar] User data received:', user);
        
        // Get first_name and last_name from user_metadata or directly
        let firstName = '';
        let lastName = '';
        
        if (user.user_metadata) {
            firstName = user.user_metadata.first_name || '';
            lastName = user.user_metadata.last_name || '';
        } else if (user.first_name) {
            firstName = user.first_name;
            lastName = user.last_name || '';
        }
        
        console.log('[loadUserInfoInNavbar] Extracted name:', { firstName, lastName });
        
        // Build display name
        let displayName = '';
        if (firstName && lastName) {
            displayName = `${firstName} ${lastName}`;
        } else if (firstName) {
            displayName = firstName;
        } else if (user.email) {
            // Fallback to email if no name
            displayName = user.email.split('@')[0];
        } else {
            displayName = 'Reader';
        }
        
        console.log('[loadUserInfoInNavbar] Setting display name:', displayName);
        
        // Update the right-side user info nav if elements exist
        if (userInfoNav && userNameDisplay) {
            userNameDisplay.textContent = displayName;
            userInfoNav.style.display = 'block';
            console.log('[loadUserInfoInNavbar] Name set in user info nav');
        }
        
        console.log('[loadUserInfoInNavbar] User info displayed successfully');
    })
    .catch(err => {
        console.error('[loadUserInfoInNavbar] Error loading user info:', err);
    });
}

// Initialize HTMX configuration
document.addEventListener('DOMContentLoaded', function() {
    console.log('[app.js] DOMContentLoaded fired');
    loadAliceReaderLanguagePreference();
    watchAliceReaderLanguageMutations();
    
    // Sync cookie from sessionStorage on page load (for server-side navigation)
    // Only sync once per page load
    if (!window.cookieSyncDone) {
        window.cookieSyncDone = true;
        ensureCookieSync();
    }
    
    // Configure HTMX
    htmx.config.globalViewTransitions = true;
    
    // Attach logout handler to all logout links (for reader pages)
    const logoutLinks = document.querySelectorAll('#logout-link-reader, a[onclick*="logout"]');
    console.log('[app.js] Found logout links:', logoutLinks.length);
    
    logoutLinks.forEach(function(link) {
        console.log('[app.js] Attaching logout handler to:', link.id || link.textContent);
        link.addEventListener('click', function(e) {
            e.preventDefault();
            console.log('[app.js] Logout link clicked');
            if (window.logout) {
                console.log('[app.js] Calling window.logout()');
                window.logout();
            } else {
                console.error('[app.js] window.logout is not defined!');
            }
        });
        // Remove onclick attribute if present
        if (link.hasAttribute('onclick')) {
            link.removeAttribute('onclick');
        }
    });
    
    // Load user info in navbar (for reader pages only, NOT consultant pages)
    // Only load if we're on a reader page (not landing page, not consultant pages)
    const path = window.location.pathname;
    const isConsultantPage = path.startsWith('/consultant');
    const isLandingPage = path === '/' || path === '/login' || path === '/register';
    
    if (!isLandingPage && !isConsultantPage && !window.isConsultantDashboard) {
        console.log('[app.js] Calling loadUserInfoInNavbar() for reader page');
        loadUserInfoInNavbar();
    } else {
        // On landing/login/register/consultant pages, ensure user info is hidden (consultant handles its own)
        if (isConsultantPage || window.isConsultantDashboard) {
            console.log('[app.js] Skipping loadUserInfoInNavbar() - consultant page, consultant handles its own');
        } else {
            const userInfoNav = document.getElementById('user-info-nav');
            if (userInfoNav) {
                userInfoNav.style.display = 'none';
            }
        }
    }
    
    // Add auth token to all HTMX requests
    document.body.addEventListener('htmx:configRequest', function(event) {
        const token = getAuthToken();
        if (token) {
            event.detail.headers['Authorization'] = 'Bearer ' + token;
        }
    });
    
    // Handle 401 errors (unauthorized)
    document.body.addEventListener('htmx:responseError', function(event) {
        if (event.detail.xhr.status === 401) {
            removeAuthToken();
            disconnectSSE();
            window.location.href = '/reader/login';
        }
    });

    // Connect to SSE if authenticated (but not on consultant dashboard - it handles its own SSE)
    // Check flag both ways to be safe
    if (isAuthenticated() && !window.isConsultantDashboard && window.location.pathname.indexOf('/consultant') === -1) {
        connectSSE();
    }
});

// Disconnect SSE on page unload
window.addEventListener('beforeunload', function() {
    disconnectSSE();
});
