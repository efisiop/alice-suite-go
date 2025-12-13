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
        const data = JSON.parse(event.data);
        handleSSEEvent(data);
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
    
    // Hide name in navbar brand as well
    const userNameBrand = document.getElementById('user-name-brand');
    if (userNameBrand) {
        userNameBrand.style.display = 'none';
        userNameBrand.textContent = '';
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
        
        // Update the name in the navbar brand (left side, after "Alice Suite") - always do this
        const userNameBrand = document.getElementById('user-name-brand');
        console.log('[loadUserInfoInNavbar] Looking for user-name-brand element:', !!userNameBrand);
        if (userNameBrand) {
            userNameBrand.textContent = displayName;
            userNameBrand.style.display = 'inline';
            console.log('[loadUserInfoInNavbar] Name set in navbar brand:', displayName);
        } else {
            console.error('[loadUserInfoInNavbar] user-name-brand element NOT FOUND!');
        }
        
        // Also update the right-side user info nav if elements exist
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

