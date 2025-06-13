// API Configuration
const API_BASE_URL = 'http://localhost:5000'; // Ajuste para a porta da sua API
let currentUser = null;
let authToken = null;
let followingUsersCache = new Set(); // Cache dos IDs dos usuários seguidos

// Dark Mode Functions
function createThemeToggle() {
    // Check if button already exists
    if (document.getElementById('theme-toggle')) {
        return;
    }
    
    console.log('Creating theme toggle button...');
    
    const button = document.createElement('button');
    button.id = 'theme-toggle';
    button.onclick = toggleTheme;
    button.style.cssText = `
        position: fixed;
        top: 10px;
        right: 10px;
        z-index: 99999;
        background: #ffffff;
        border: 1px solid #333333;
        border-radius: 20px;
        padding: 4px 8px;
        cursor: pointer;
        display: flex !important;
        align-items: center;
        gap: 4px;
        font-size: 0.7rem;
        font-weight: 500;
        color: #333333;
        box-shadow: 0 2px 6px rgba(0,0,0,0.15);
        transition: all 0.3s ease;
        font-family: inherit;
    `;
    
    button.innerHTML = `
        <span class="icon" style="font-size: 0.8rem;">🌙</span>
        <span class="text">Dark Mode</span>
    `;
    
    document.body.appendChild(button);
    console.log('Theme toggle button created and added to page');
}

function initializeTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    console.log('Initializing theme:', savedTheme);
    
    // Ensure button exists
    createThemeToggle();
    
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeToggle(savedTheme);
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    console.log('Toggling theme from', currentTheme, 'to', newTheme);
    
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeToggle(newTheme);
    
    // Aplicar mudanças visuais imediatas
    if (newTheme === 'dark') {
        document.body.style.background = 'linear-gradient(135deg, #7c3aed 0%, #a855f7 100%)';
        document.body.style.color = '#ffffff';
    } else {
        document.body.style.background = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)';
        document.body.style.color = '#333333';
    }
    
    // Animação suave
    document.body.style.transition = 'all 0.3s ease';
    setTimeout(() => {
        document.body.style.transition = '';
    }, 300);
}

function updateThemeToggle(theme) {
    const toggle = document.getElementById('theme-toggle') || document.querySelector('.theme-toggle');
    if (!toggle) {
        console.error('Theme toggle button not found!');
        return;
    }
    
    const icon = toggle.querySelector('.icon');
    const text = toggle.querySelector('.text');
    
    if (!icon || !text) {
        console.error('Theme toggle elements not found!');
        return;
    }
    
    console.log('Updating theme toggle for theme:', theme);
    
    if (theme === 'dark') {
        icon.textContent = '☀️';
        icon.style.fontSize = '0.8rem';
        text.textContent = 'Light Mode';
        toggle.setAttribute('title', 'Mudar para modo claro');
        toggle.style.background = '#2d2d2d';
        toggle.style.color = '#ffffff';
        toggle.style.borderColor = '#ffffff';
    } else {
        icon.textContent = '🌙';
        icon.style.fontSize = '0.8rem';
        text.textContent = 'Dark Mode';
        toggle.setAttribute('title', 'Mudar para modo escuro');
        toggle.style.background = '#ffffff';
        toggle.style.color = '#333333';
        toggle.style.borderColor = '#333333';
    }
}

// Utility Functions
function showLoader() {
    document.getElementById('loader').style.display = 'flex';
}

function hideLoader() {
    document.getElementById('loader').style.display = 'none';
}

function showToast(message, type = 'info') {
    const toastContainer = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <div>
            <strong>${type === 'success' ? 'Sucesso!' : type === 'error' ? 'Erro!' : 'Info!'}</strong>
            <p>${message}</p>
        </div>
    `;
    
    toastContainer.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 5000);
}

function showModal(title, message, onConfirm) {
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modal-title');
    const modalMessage = document.getElementById('modal-message');
    const modalConfirm = document.getElementById('modal-confirm');
    
    modalTitle.textContent = title;
    modalMessage.textContent = message;
    modal.style.display = 'flex';
    
    modalConfirm.onclick = () => {
        modal.style.display = 'none';
        if (onConfirm) onConfirm();
    };
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
}

// API Functions
async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    
    if (authToken) {
        defaultOptions.headers['Authorization'] = `Bearer ${authToken}`;
    }
    
    const config = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...options.headers,
        },
    };
    
    try {
        const response = await fetch(url, config);
        
        if (!response.ok) {
            let errorMessage = 'Erro na requisição';
            try {
                const errorText = await response.text();
                if (errorText) {
                    try {
                        const errorJson = JSON.parse(errorText);
                        errorMessage = errorJson.erro || errorText;
                    } catch (e) {
                        errorMessage = errorText;
                    }
                } else {
                    errorMessage = `Erro ${response.status}`;
                }
            } catch (e) {
                errorMessage = `Erro ${response.status}`;
            }
            throw new Error(errorMessage);
        }
        
        // Handle empty responses (204 No Content)
        if (response.status === 204) {
            return null;
        }
        
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
            const text = await response.text();
            return text ? JSON.parse(text) : null;
        } else {
            return await response.text();
        }
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// Authentication Functions
function setAuthToken(token) {
    authToken = token;
    localStorage.setItem('authToken', token);
}

function getAuthToken() {
    return localStorage.getItem('authToken');
}

function clearAuthToken() {
    authToken = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
}

// Decode JWT token to get user ID
function decodeToken(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        return JSON.parse(jsonPayload);
    } catch (error) {
        console.error('Error decoding token:', error);
        return null;
    }
}

// Auth Form Functions
function showLogin() {
    document.getElementById('login-form').style.display = 'block';
    document.getElementById('register-form').style.display = 'none';
}

function showRegister() {
    document.getElementById('login-form').style.display = 'none';
    document.getElementById('register-form').style.display = 'block';
}

// Login Function
async function login(event) {
    event.preventDefault();
    showLoader();
    
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    
    try {
        const response = await apiRequest('/login', {
            method: 'POST',
            body: JSON.stringify({
                email: email,
                senha: password
            })
        });
        
        // Response should be the token as plain text
        const token = typeof response === 'string' ? response : response.token;
        setAuthToken(token);
        
        // Get user info
        const tokenData = decodeToken(token);
        if (tokenData && tokenData.usuarioId) {
            const userInfo = await apiRequest(`/usuarios/${tokenData.usuarioId}`);
            currentUser = userInfo;
            localStorage.setItem('currentUser', JSON.stringify(userInfo));
            
            // Carregar cache de usuários seguidos
            await loadFollowingUsersCache();
        }
        
        showToast('Login realizado com sucesso!', 'success');
        showApp();
    } catch (error) {
        showToast(error.message, 'error');
    } finally {
        hideLoader();
    }
}

// Register Function
async function register(event) {
    event.preventDefault();
    showLoader();
    
    const name = document.getElementById('register-name').value;
    const nick = document.getElementById('register-nick').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    
    try {
        await apiRequest('/usuarios', {
            method: 'POST',
            body: JSON.stringify({
                nome: name,
                nick: nick,
                email: email,
                senha: password
            })
        });
        
        showToast('Usuário criado com sucesso! Faça login para continuar.', 'success');
        showLogin();
        
        // Clear form
        document.getElementById('register-name').value = '';
        document.getElementById('register-nick').value = '';
        document.getElementById('register-email').value = '';
        document.getElementById('register-password').value = '';
    } catch (error) {
        showToast(error.message, 'error');
    } finally {
        hideLoader();
    }
}

// App Functions
function showApp() {
    document.getElementById('auth-section').style.display = 'none';
    document.getElementById('app-content').style.display = 'block';
    document.getElementById('navbar').style.display = 'block';
    
    showHome();
    loadDashboardStats();
}

function showHome() {
    hideAllSections();
    document.getElementById('home-section').style.display = 'block';
    setActiveNavButton(0);
}

function showProfile() {
    hideAllSections();
    document.getElementById('profile-section').style.display = 'block';
    setActiveNavButton(1);
    loadProfile();
}

function showUsers() {
    hideAllSections();
    document.getElementById('users-section').style.display = 'block';
    setActiveNavButton(2);
    loadUsers();
}

function hideAllSections() {
    const sections = document.querySelectorAll('.section');
    sections.forEach(section => section.style.display = 'none');
}

function setActiveNavButton(index) {
    const buttons = document.querySelectorAll('.nav-btn');
    buttons.forEach((btn, i) => {
        if (i === index) {
            btn.style.background = '#667eea';
            btn.style.color = 'white';
        } else if (!btn.classList.contains('logout-btn')) {
            btn.style.background = 'none';
            btn.style.color = '#333';
        }
    });
}

// Dashboard Functions
async function loadDashboardStats() {
    try {
        const users = await apiRequest('/usuarios');
        document.getElementById('total-users').textContent = users.length || 0;
        
        // Buscar estatísticas do usuário atual
        if (currentUser && currentUser.id) {
            const stats = await apiRequest(`/usuarios/${currentUser.id}/estatisticas`);
            document.getElementById('following-count').textContent = stats.seguindo || 0;
            document.getElementById('followers-count').textContent = stats.seguidores || 0;
        } else {
            // Tentar carregar do localStorage
            const storedUser = localStorage.getItem('currentUser');
            if (storedUser) {
                currentUser = JSON.parse(storedUser);
                if (currentUser && currentUser.id) {
                    const stats = await apiRequest(`/usuarios/${currentUser.id}/estatisticas`);
                    document.getElementById('following-count').textContent = stats.seguindo || 0;
                    document.getElementById('followers-count').textContent = stats.seguidores || 0;
                    return;
                }
            }
            document.getElementById('following-count').textContent = '0';
            document.getElementById('followers-count').textContent = '0';
        }
    } catch (error) {
        console.error('Error loading dashboard stats:', error);
        // Em caso de erro, manter valores padrão
        document.getElementById('following-count').textContent = '0';
        document.getElementById('followers-count').textContent = '0';
    }
}

// Profile Functions
async function loadProfileWithPhoto() {
    try {
        const userId = getCurrentUserId();
        const user = await apiRequest(`/usuarios/${userId}`);
        
        document.getElementById('profile-name').textContent = user.nome;
        document.getElementById('profile-email').textContent = user.email;
        document.getElementById('profile-nick').textContent = `@${user.nick}`;
        
        // Carregar foto de perfil se existir
        updateProfilePhoto(user.foto_perfil);
        
        // Update form fields
        document.getElementById('edit-name').value = user.nome;
        document.getElementById('edit-nick').value = user.nick;
        document.getElementById('edit-email').value = user.email;
        
        currentUser = user;
        localStorage.setItem('currentUser', JSON.stringify(user));
        
    } catch (error) {
        console.error('Erro ao carregar perfil:', error);
        showToast('Erro ao carregar dados do perfil', 'error');
    }
}

function editProfile() {
    document.getElementById('edit-profile-form').style.display = 'block';
    
    // Fill form with current data
    document.getElementById('edit-name').value = currentUser.nome || '';
    document.getElementById('edit-nick').value = currentUser.nick || '';
    document.getElementById('edit-email').value = currentUser.email || '';
}

function cancelEdit() {
    document.getElementById('edit-profile-form').style.display = 'none';
}

async function updateProfile(event) {
    event.preventDefault();
    showLoader();
    
    const name = document.getElementById('edit-name').value;
    const nick = document.getElementById('edit-nick').value;
    const email = document.getElementById('edit-email').value;
    
    try {
        await apiRequest(`/usuarios/${currentUser.id}`, {
            method: 'PUT',
            body: JSON.stringify({
                nome: name,
                nick: nick,
                email: email
            })
        });
        
        // Update current user data
        currentUser.nome = name;
        currentUser.nick = nick;
        currentUser.email = email;
        localStorage.setItem('currentUser', JSON.stringify(currentUser));
        
        showToast('Perfil atualizado com sucesso!', 'success');
        loadProfile();
        cancelEdit();
    } catch (error) {
        showToast(error.message, 'error');
    } finally {
        hideLoader();
    }
}

function deleteAccount() {
    showModal(
        'Excluir Conta',
        'Tem certeza que deseja excluir sua conta? Esta ação não pode ser desfeita.',
        async () => {
            showLoader();
            try {
                await apiRequest(`/usuarios/${currentUser.id}`, {
                    method: 'DELETE'
                });
                
                showToast('Conta excluída com sucesso!', 'success');
                logout();
            } catch (error) {
                showToast(error.message, 'error');
            } finally {
                hideLoader();
            }
        }
    );
}

// Cache Functions
async function loadFollowingUsersCache() {
    console.log('🔄 Carregando cache de usuários seguidos...');
    console.log('currentUser:', currentUser);
    
    if (!currentUser || !currentUser.id) {
        const storedUser = localStorage.getItem('currentUser');
        if (storedUser) {
            currentUser = JSON.parse(storedUser);
            console.log('📦 currentUser carregado do localStorage:', currentUser);
        }
    }
    
    if (currentUser && currentUser.id) {
        try {
            console.log(`🌐 Fazendo requisição para /usuarios/${currentUser.id}/seguindo`);
            const followingUsers = await apiRequest(`/usuarios/${currentUser.id}/seguindo`);
            followingUsersCache.clear();
            
            // Tratar caso quando API retorna null (usuário não segue ninguém)
            if (followingUsers && Array.isArray(followingUsers)) {
                followingUsers.forEach(user => followingUsersCache.add(user.id));
                console.log('✅ Cache carregado com sucesso:', Array.from(followingUsersCache));
            } else {
                console.log('✅ Cache carregado - usuário não segue ninguém');
            }
        } catch (error) {
            console.error('❌ Erro ao carregar cache de usuários seguidos:', error);
        }
    } else {
        console.log('⚠️ currentUser não definido ou sem ID');
    }
}

function updateFollowingCache(userId, isFollowing) {
    if (isFollowing) {
        followingUsersCache.add(userId);
    } else {
        followingUsersCache.delete(userId);
    }
}

// Users Functions
async function loadUsers(searchTerm = '') {
    console.log('👥 Carregando usuários...');
    console.log('Cache atual:', Array.from(followingUsersCache));
    
    showLoader();
    try {
        let endpoint = '/usuarios';
        if (searchTerm) {
            endpoint += `?usuario=${encodeURIComponent(searchTerm)}`;
        }
        
        const users = await apiRequest(endpoint);
        
        // Carregar cache de usuários seguidos se estiver vazio
        if (followingUsersCache.size === 0) {
            console.log('🔄 Cache vazio, carregando...');
            await loadFollowingUsersCache();
        } else {
            console.log('✅ Usando cache existente:', Array.from(followingUsersCache));
        }
        
        displayUsers(users);
    } catch (error) {
        showToast('Erro ao carregar usuários: ' + error.message, 'error');
        document.getElementById('users-list').innerHTML = '<p>Erro ao carregar usuários</p>';
    } finally {
        hideLoader();
    }
}

function displayUsers(users) {
    console.log('🎨 Renderizando usuários...');
    console.log('Cache para renderização:', Array.from(followingUsersCache));
    
    const usersList = document.getElementById('users-list');
    
    if (!users || users.length === 0) {
        usersList.innerHTML = '<p style="color: white; text-align: center;">Nenhum usuário encontrado</p>';
        return;
    }
    
    usersList.innerHTML = users.map(user => {
        const isFollowing = followingUsersCache.has(user.id);
        console.log(`👤 Usuário ${user.id} (${user.nome}): seguindo = ${isFollowing}`);
        
        return `
        <div class="user-card">
            <div class="user-avatar">
                <i class="fas fa-user-circle"></i>
            </div>
            <div class="user-info">
                <h3>${user.nome}</h3>
                <p>@${user.nick}</p>
                <p>${user.email}</p>
            </div>
            ${user.id !== currentUser?.id ? `
                <div class="user-actions" data-user-id="${user.id}">
                    <button class="btn btn-sm btn-follow" data-user-id="${user.id}" data-action="follow" style="display: ${isFollowing ? 'none' : 'inline-flex'};">
                        <i class="fas fa-user-plus"></i> Seguir
                    </button>
                    <button class="btn btn-sm btn-unfollow" data-user-id="${user.id}" data-action="unfollow" style="display: ${isFollowing ? 'inline-flex' : 'none'};">
                        <i class="fas fa-user-minus"></i> Deixar de seguir
                    </button>
                </div>
            ` : '<p style="color: #667eea; font-weight: bold;">Você</p>'}
        </div>
        `;
    }).join('');
    
    // Adicionar event listeners para os botões
    addUserActionListeners();
}

// Adicionar event listeners para botões de seguir/deixar de seguir
function addUserActionListeners() {
    // Remover listeners antigos para evitar duplicação
    document.querySelectorAll('[data-action="follow"], [data-action="unfollow"]').forEach(button => {
        button.removeEventListener('click', handleUserAction);
        button.addEventListener('click', handleUserAction);
    });
}

// Handler para ações de seguir/deixar de seguir
async function handleUserAction(event) {
    event.preventDefault();
    const button = event.target.closest('button');
    const userId = parseInt(button.dataset.userId);
    const action = button.dataset.action;
    
    if (action === 'follow') {
        await followUser(userId, event);
    } else if (action === 'unfollow') {
        await unfollowUser(userId, event);
    }
}

// Função para atualizar os botões de seguir/deixar de seguir
function updateUserButtons(userId, action) {
    const userActions = document.querySelector(`[data-user-id="${userId}"]`);
    if (userActions) {
        const followBtn = userActions.querySelector('.btn-follow');
        const unfollowBtn = userActions.querySelector('.btn-unfollow');
        
        if (followBtn && unfollowBtn) {
            if (action === 'followed') {
                followBtn.style.display = 'none';
                unfollowBtn.style.display = 'inline-flex';
            } else if (action === 'unfollowed') {
                followBtn.style.display = 'inline-flex';
                unfollowBtn.style.display = 'none';
            }
        }
    }
}

async function searchUsers() {
    const searchTerm = document.getElementById('user-search').value.trim();
    await loadUsers(searchTerm);
}

// Follow/Unfollow Functions
async function followUser(userId, event) {
    showLoader();
    try {
        await apiRequest(`/usuarios/${userId}/seguir`, {
            method: 'POST'
        });
        
        showToast('Usuário seguido com sucesso!', 'success');
        
        // Atualizar cache
        updateFollowingCache(userId, true);
        
        // Update UI - usar abordagem robusta
        updateUserButtons(userId, 'followed');
        
        // Atualizar estatísticas na dashboard
        loadDashboardStats();
        
    } catch (error) {
        console.error('Erro ao seguir usuário:', error);
        
        // Tratamento específico de erros
        if (error.message.includes('não é possível seguir você mesmo')) {
            showToast('Você não pode seguir a si mesmo!', 'error');
        } else if (error.message.includes('usuário não encontrado')) {
            showToast('Usuário não encontrado!', 'error');
        } else {
            showToast(`Erro ao seguir usuário: ${error.message}`, 'error');
        }
    } finally {
        hideLoader();
    }
}

async function unfollowUser(userId, event) {
    showLoader();
    try {
        await apiRequest(`/usuarios/${userId}/parar-de-seguir`, {
            method: 'POST'
        });
        
        showToast('Você parou de seguir este usuário', 'success');
        
        // Atualizar cache
        updateFollowingCache(userId, false);
        
        // Update UI - usar abordagem robusta
        updateUserButtons(userId, 'unfollowed');
        
        // Atualizar estatísticas na dashboard
        loadDashboardStats();
        
    } catch (error) {
        console.error('Erro ao deixar de seguir usuário:', error);
        
        // Tratamento específico de erros
        if (error.message.includes('não é possível parar de seguir voce mesmo')) {
            showToast('Você não pode deixar de seguir a si mesmo!', 'error');
        } else {
            showToast(`Erro ao deixar de seguir usuário: ${error.message}`, 'error');
        }
    } finally {
        hideLoader();
    }
}

// Logout Function
function logout() {
    clearAuthToken();
    currentUser = null;
    followingUsersCache.clear(); // Limpar cache
    
    document.getElementById('auth-section').style.display = 'block';
    document.getElementById('app-content').style.display = 'none';
    document.getElementById('navbar').style.display = 'none';
    
    // Clear forms
    document.getElementById('login-email').value = '';
    document.getElementById('login-password').value = '';
    
    showLogin();
    showToast('Logout realizado com sucesso!', 'success');
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function() {
    hideLoader();
    
    // Check if user is already logged in
    const storedToken = getAuthToken();
    const storedUser = localStorage.getItem('currentUser');
    
    if (storedToken && storedUser) {
        authToken = storedToken;
        currentUser = JSON.parse(storedUser);
        // Carregar cache de usuários seguidos
        loadFollowingUsersCache();
        showApp();
    } else {
        showLogin();
    }
    
    // Search on Enter key
    document.getElementById('user-search').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchUsers();
        }
    });
    
    // Close modal when clicking outside
    document.getElementById('modal').addEventListener('click', function(e) {
        if (e.target === this) {
            closeModal();
        }
    });
});

// Error handling for uncaught errors
window.addEventListener('error', function(e) {
    console.error('Global error:', e.error);
    showToast('Ocorreu um erro inesperado', 'error');
});

// Handle network errors
window.addEventListener('online', function() {
    showToast('Conexão restaurada', 'success');
});

window.addEventListener('offline', function() {
    showToast('Conexão perdida', 'warning');
});

// === POSTS FUNCTIONS ===

// Show feed section
function showFeed() {
    hideAllSections();
    document.getElementById('feed-section').style.display = 'block';
    loadPosts();
}

// Show create post form
function showCreatePost() {
    document.getElementById('create-post-form').style.display = 'block';
    document.getElementById('post-title').focus();
}

// Cancel create post
function cancelCreatePost() {
    document.getElementById('create-post-form').style.display = 'none';
    document.getElementById('post-title').value = '';
    document.getElementById('post-content').value = '';
}

// Create new post
async function createPost(event) {
    event.preventDefault();
    
    const title = document.getElementById('post-title').value.trim();
    const content = document.getElementById('post-content').value.trim();
    
    if (!title || !content) {
        showToast('Por favor, preencha todos os campos', 'error');
        return;
    }
    
    try {
        await apiRequest('/publicacoes', {
            method: 'POST',
            body: JSON.stringify({
                titulo: title,
                conteudo: content
            })
        });
        
        showToast('Publicação criada com sucesso!', 'success');
        cancelCreatePost();
        loadPosts();
    } catch (error) {
        showToast(`Erro ao criar publicação: ${error.message}`, 'error');
    }
}

// Load posts
async function loadPosts() {
    const postsContainer = document.getElementById('posts-list');
    
    try {
        // Show loading
        postsContainer.innerHTML = '<div class="loading">Carregando publicações...</div>';
        
        const posts = await apiRequest('/publicacoes/feed');
        displayPosts(posts);
    } catch (error) {
        console.error('Error loading posts:', error);
        postsContainer.innerHTML = '<div class="no-posts"><i class="fas fa-exclamation-triangle"></i><p>Erro ao carregar publicações</p></div>';
    }
}

// Display posts
function displayPosts(posts) {
    const postsContainer = document.getElementById('posts-list');
    
    if (!posts || posts.length === 0) {
        postsContainer.innerHTML = '<div class="no-posts"><i class="fas fa-newspaper"></i><p>Nenhuma publicação encontrada</p></div>';
        return;
    }
    
    const currentUserId = getCurrentUserId();
    
    postsContainer.innerHTML = posts.map(post => {
        const createdAt = new Date(post.criadoem).toLocaleString('pt-BR');
        const isOwner = post.autorId === currentUserId;
        
        return `
            <div class="post-card" data-post-id="${post.id}">
                <div class="post-header">
                    <div class="post-author">
                        <div class="post-author-avatar">
                            ${post.autorNick.charAt(0).toUpperCase()}
                        </div>
                        <div class="post-author-info">
                            <h4>${post.autorNome}</h4>
                            <p>@${post.autorNick}</p>
                        </div>
                    </div>
                    <div class="post-time">${createdAt}</div>
                </div>
                
                <div class="post-content">
                    <h3>${post.titulo}</h3>
                    <p>${post.conteudo}</p>
                </div>
                
                <div class="post-actions">
                    <button class="post-action-btn" onclick="likePost(${post.id})">
                        <i class="fas fa-heart"></i>
                        <span>${post.curtidas || 0}</span>
                    </button>
                    <button class="post-action-btn" onclick="toggleComments(${post.id})">
                        <i class="fas fa-comment"></i>
                        <span>${post.comentarios || 0}</span>
                    </button>
                    ${isOwner ? `
                        <button class="post-action-btn edit-btn" onclick="editPost(${post.id})">
                            <i class="fas fa-edit"></i>
                            Editar
                        </button>
                        <button class="post-action-btn delete-btn" onclick="deletePost(${post.id})">
                            <i class="fas fa-trash"></i>
                            Excluir
                        </button>
                    ` : ''}
                </div>
                
                <!-- Comments Section -->
                <div class="comments-section" id="comments-${post.id}" style="display: none;">
                    <div class="comment-form">
                        <textarea id="comment-input-${post.id}" placeholder="💬 Compartilhe sua opinião sobre esta publicação..." maxlength="500"></textarea>
                        <button onclick="addComment(${post.id})" class="btn btn-sm btn-primary">
                            <i class="fas fa-paper-plane"></i> Comentar
                        </button>
                    </div>
                    <div class="comments-list" id="comments-list-${post.id}">
                        <!-- Comments will be loaded here -->
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

// Like post
async function likePost(postId) {
    try {
        await apiRequest(`/publicacoes/${postId}/curtir`, {
            method: 'POST'
        });
        
        loadPosts(); // Reload to update like count
    } catch (error) {
        showToast('Erro ao curtir publicação', 'error');
    }
}

// Delete post
async function deletePost(postId) {
    if (!confirm('Tem certeza que deseja excluir esta publicação?')) {
        return;
    }
    
    try {
        await apiRequest(`/publicacoes/${postId}`, {
            method: 'DELETE'
        });
        
        showToast('Publicação excluída com sucesso!', 'success');
        loadPosts();
    } catch (error) {
        showToast('Erro ao excluir publicação', 'error');
    }
}

// Edit post (placeholder for now)
function editPost(postId) {
    showToast('Funcionalidade de edição em desenvolvimento', 'info');
}

// Helper to get current user ID from token
function getCurrentUserId() {
    const token = getAuthToken();
    if (!token) return null;
    
    try {
        const payload = decodeToken(token);
        return payload.usuarioId;
    } catch (error) {
        console.error('Error parsing token:', error);
        return null;
    }
}

// Comments Functions
async function toggleComments(postId) {
    const commentsSection = document.getElementById(`comments-${postId}`);
    
    if (commentsSection.style.display === 'none') {
        commentsSection.style.display = 'block';
        await loadComments(postId);
    } else {
        commentsSection.style.display = 'none';
    }
}

async function loadComments(postId) {
    const commentsList = document.getElementById(`comments-list-${postId}`);
    
    try {
        commentsList.innerHTML = '<div class="loading">Carregando comentários...</div>';
        
        const comments = await apiRequest(`/publicacoes/${postId}/comentarios`);
        displayComments(postId, comments);
    } catch (error) {
        console.error('Error loading comments:', error);
        commentsList.innerHTML = '<div class="error">Erro ao carregar comentários</div>';
    }
}

function displayComments(postId, comments) {
    const commentsList = document.getElementById(`comments-list-${postId}`);
    const currentUserId = getCurrentUserId();
    
    if (!comments || comments.length === 0) {
        commentsList.innerHTML = '<div class="no-comments">Nenhum comentário ainda</div>';
        return;
    }
    
    commentsList.innerHTML = comments.map(comment => {
        const createdAt = new Date(comment.criadoem).toLocaleString('pt-BR');
        const isOwner = comment.autorId === currentUserId;
        
        return `
            <div class="comment" data-comment-id="${comment.id}">
                <div class="comment-header">
                    <div class="comment-author">
                        <div class="comment-author-avatar">
                            ${comment.autorNick.charAt(0).toUpperCase()}
                        </div>
                        <div class="comment-author-info">
                            <strong>${comment.autorNome}</strong>
                            <span>@${comment.autorNick}</span>
                            <span class="comment-time">${createdAt}</span>
                        </div>
                    </div>
                    ${isOwner ? `
                        <div class="comment-actions">
                            <button class="comment-action-btn" onclick="editComment(${comment.id})">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="comment-action-btn" onclick="deleteComment(${comment.id}, ${postId})">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    ` : ''}
                </div>
                <div class="comment-content" id="comment-content-${comment.id}">
                    <p>${comment.conteudo}</p>
                </div>
                <div class="comment-reactions" id="reactions-${comment.id}">
                    <div class="reaction-buttons">
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'like')">👍 <span id="like-count-${comment.id}">0</span></button>
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'love')">❤️ <span id="love-count-${comment.id}">0</span></button>
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'laugh')">😂 <span id="laugh-count-${comment.id}">0</span></button>
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'wow')">😮 <span id="wow-count-${comment.id}">0</span></button>
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'sad')">😢 <span id="sad-count-${comment.id}">0</span></button>
                        <button class="reaction-btn" onclick="toggleReaction(${comment.id}, 'angry')">😡 <span id="angry-count-${comment.id}">0</span></button>
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

async function addComment(postId) {
    const commentInput = document.getElementById(`comment-input-${postId}`);
    const content = commentInput.value.trim();
    
    if (!content) {
        showToast('Por favor, escreva um comentário', 'error');
        return;
    }
    
    try {
        await apiRequest(`/publicacoes/${postId}/comentarios`, {
            method: 'POST',
            body: JSON.stringify({
                conteudo: content
            })
        });
        
        commentInput.value = '';
        showToast('Comentário adicionado com sucesso!', 'success');
        
        // Reload comments and update post
        await loadComments(postId);
        loadPosts(); // Reload to update comment count
    } catch (error) {
        showToast(`Erro ao adicionar comentário: ${error.message}`, 'error');
    }
}

async function editComment(commentId) {
    const commentContent = document.getElementById(`comment-content-${commentId}`);
    const currentContent = commentContent.querySelector('p').textContent;
    
    const newContent = prompt('Editar comentário:', currentContent);
    if (newContent === null || newContent.trim() === '') return;
    
    try {
        await apiRequest(`/comentarios/${commentId}`, {
            method: 'PUT',
            body: JSON.stringify({
                conteudo: newContent.trim()
            })
        });
        
        commentContent.querySelector('p').textContent = newContent.trim();
        showToast('Comentário atualizado com sucesso!', 'success');
    } catch (error) {
        showToast(`Erro ao atualizar comentário: ${error.message}`, 'error');
    }
}

async function deleteComment(commentId, postId) {
    if (!confirm('Tem certeza que deseja excluir este comentário?')) {
        return;
    }
    
    try {
        await apiRequest(`/comentarios/${commentId}`, {
            method: 'DELETE'
        });
        
        showToast('Comentário excluído com sucesso!', 'success');
        
        // Reload comments and update post
        await loadComments(postId);
        loadPosts(); // Reload to update comment count
    } catch (error) {
        showToast(`Erro ao excluir comentário: ${error.message}`, 'error');
    }
}

// Reaction Functions (Frontend only - simplified version)
const reactionCounts = {}; // Store reaction counts in memory
const userReactions = {}; // Store user's reactions in memory

function toggleReaction(commentId, reactionType) {
    // Initialize if not exists
    if (!reactionCounts[commentId]) {
        reactionCounts[commentId] = {
            like: 0, love: 0, laugh: 0, wow: 0, sad: 0, angry: 0
        };
    }
    if (!userReactions[commentId]) {
        userReactions[commentId] = {};
    }
    
    const countSpan = document.getElementById(`${reactionType}-count-${commentId}`);
    const reactionBtn = countSpan.parentElement;
    
    // Check if user already reacted with this type
    const hasReacted = userReactions[commentId][reactionType];
    
    if (hasReacted) {
        // Remove reaction
        reactionCounts[commentId][reactionType]--;
        userReactions[commentId][reactionType] = false;
        reactionBtn.classList.remove('active');
        showToast('Reação removida!', 'success');
    } else {
        // Add reaction
        reactionCounts[commentId][reactionType]++;
        userReactions[commentId][reactionType] = true;
        reactionBtn.classList.add('active');
        showToast(`Você reagiu com ${getReactionEmoji(reactionType)}!`, 'success');
    }
    
    // Update count display with animation
    countSpan.textContent = reactionCounts[commentId][reactionType];
    countSpan.classList.add('animate');
    setTimeout(() => countSpan.classList.remove('animate'), 300);
}

function getReactionEmoji(reactionType) {
    const emojis = {
        'like': '👍',
        'love': '❤️',
        'laugh': '😂',
        'wow': '😮',
        'sad': '😢',
        'angry': '😡'
    };
    return emojis[reactionType] || '👍';
}

// Initialize theme when page loads
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing theme...');
    initializeTheme();
    
    // Check for existing auth token
    const token = getAuthToken();
    if (token) {
        authToken = token;
        const savedUser = localStorage.getItem('currentUser');
        if (savedUser) {
            currentUser = JSON.parse(savedUser);
            showApp();
        }
    }
});

// Also initialize on window load as backup
window.addEventListener('load', function() {
    console.log('Window loaded, ensuring theme is initialized...');
    if (!document.documentElement.getAttribute('data-theme')) {
        initializeTheme();
    }
});

// Photo Upload Functions
function openPhotoUpload() {
    document.getElementById('photo-upload-form').style.display = 'block';
    document.getElementById('edit-profile-form').style.display = 'none';
}

function cancelPhotoUpload() {
    document.getElementById('photo-upload-form').style.display = 'none';
    document.getElementById('photo-input').value = '';
    document.getElementById('photo-preview').style.display = 'none';
}

function previewPhoto(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Verificar tamanho do arquivo (10MB)
    if (file.size > 10 * 1024 * 1024) {
        showToast('Arquivo muito grande! Máximo 10MB permitido.', 'error');
        event.target.value = '';
        return;
    }

    // Verificar tipo de arquivo
    if (!file.type.startsWith('image/')) {
        showToast('Apenas arquivos de imagem são permitidos!', 'error');
        event.target.value = '';
        return;
    }

    const reader = new FileReader();
    reader.onload = function(e) {
        const previewContainer = document.getElementById('photo-preview');
        const previewImage = document.getElementById('preview-image');
        
        previewImage.src = e.target.result;
        previewContainer.style.display = 'block';
    };
    reader.readAsDataURL(file);
}

async function uploadPhoto(event) {
    event.preventDefault();
    
    const fileInput = document.getElementById('photo-input');
    const file = fileInput.files[0];
    
    if (!file) {
        showToast('Por favor, selecione uma foto!', 'error');
        return;
    }

    try {
        showLoader();
        
        const formData = new FormData();
        formData.append('foto', file);
        
        const userId = getCurrentUserId();
        const response = await fetch(`${API_BASE_URL}/usuarios/${userId}/foto-perfil`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`
            },
            body: formData
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.erro || 'Erro ao fazer upload da foto');
        }

        const result = await response.json();
        
        // Atualizar a foto de perfil na interface
        updateProfilePhoto(result.foto_perfil);
        
        // Fechar formulário
        cancelPhotoUpload();
        
        showToast('Foto de perfil atualizada com sucesso!', 'success');
        
    } catch (error) {
        console.error('Erro no upload:', error);
        showToast(`Erro ao fazer upload: ${error.message}`, 'error');
    } finally {
        hideLoader();
    }
}

function updateProfilePhoto(photoFileName) {
    const profilePhoto = document.getElementById('profile-photo');
    const profileIcon = document.getElementById('profile-icon');
    
    if (photoFileName) {
        profilePhoto.src = `${API_BASE_URL}/uploads/perfil/${photoFileName}`;
        profilePhoto.style.display = 'block';
        profileIcon.style.display = 'none';
        
        // Atualizar também no localStorage se necessário
        if (currentUser) {
            currentUser.foto_perfil = photoFileName;
            localStorage.setItem('currentUser', JSON.stringify(currentUser));
        }
    } else {
        profilePhoto.style.display = 'none';
        profileIcon.style.display = 'block';
    }
}

// Substituir a função loadProfile original
window.loadProfile = loadProfileWithPhoto;

// ===== HEALTH & WELLNESS FUNCTIONS =====

// Wellness Section Navigation
function showWellness() {
    hideAllSections();
    document.getElementById('wellness-section').style.display = 'block';
    setActiveNavButton(4); // Assuming wellness is the 5th nav button
    initializeWellnessData();
}

// Initialize wellness data from localStorage
function initializeWellnessData() {
    // Load mood data
    const todayMood = localStorage.getItem(`mood_${getTodayDate()}`);
    if (todayMood) {
        const moodBtn = document.querySelector(`[data-mood="${todayMood}"]`);
        if (moodBtn) {
            moodBtn.classList.add('selected');
        }
    }
    
    // Load water intake
    const waterCount = parseInt(localStorage.getItem(`water_${getTodayDate()}`) || '0');
    updateWaterDisplay(waterCount);
    
    // Load meditation timer state
    const meditationTime = parseInt(localStorage.getItem('meditation_time') || '5');
    setMeditationTime(meditationTime);
}

// Utility function to get today's date string
function getTodayDate() {
    return new Date().toISOString().split('T')[0];
}

// Mood Tracking Functions
function saveMood() {
    const selectedMood = document.querySelector('.mood-btn.selected');
    if (!selectedMood) {
        showToast('Por favor, selecione como você se sente hoje', 'error');
        return;
    }
    
    const mood = selectedMood.getAttribute('data-mood');
    localStorage.setItem(`mood_${getTodayDate()}`, mood);
    showToast('Humor registrado com sucesso! 😊', 'success');
}

// Add event listeners to mood buttons
document.addEventListener('DOMContentLoaded', function() {
    const moodButtons = document.querySelectorAll('.mood-btn');
    moodButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            // Remove selection from all buttons
            moodButtons.forEach(b => b.classList.remove('selected'));
            // Add selection to clicked button
            this.classList.add('selected');
        });
    });
});

// Water Tracking Functions
function addWater() {
    const currentCount = parseInt(localStorage.getItem(`water_${getTodayDate()}`) || '0');
    const newCount = Math.min(currentCount + 1, 8); // Max 8 glasses
    
    localStorage.setItem(`water_${getTodayDate()}`, newCount.toString());
    updateWaterDisplay(newCount);
    
    if (newCount === 8) {
        showToast('Parabéns! Você atingiu sua meta de hidratação! 💧', 'success');
    } else {
        showToast(`+1 copo registrado! (${newCount}/8)`, 'success');
    }
}

function updateWaterDisplay(count) {
    const waterCountElement = document.getElementById('water-count');
    const waterProgressElement = document.getElementById('water-progress');
    
    if (waterCountElement) {
        waterCountElement.textContent = count;
    }
    
    if (waterProgressElement) {
        const percentage = (count / 8) * 100;
        waterProgressElement.style.width = `${percentage}%`;
    }
}

// Exercise Tracking Functions
function logExercise() {
    const minutes = document.getElementById('exercise-minutes').value;
    const type = document.getElementById('exercise-type').value;
    
    if (!minutes || !type) {
        showToast('Por favor, preencha todos os campos', 'error');
        return;
    }
    
    if (minutes < 1 || minutes > 300) {
        showToast('Por favor, insira um tempo válido (1-300 minutos)', 'error');
        return;
    }
    
    // Save to localStorage
    const today = getTodayDate();
    const exercises = JSON.parse(localStorage.getItem(`exercises_${today}`) || '[]');
    exercises.push({
        type: type,
        minutes: parseInt(minutes),
        timestamp: new Date().toISOString()
    });
    localStorage.setItem(`exercises_${today}`, JSON.stringify(exercises));
    
    // Calculate total minutes today
    const totalMinutes = exercises.reduce((sum, ex) => sum + ex.minutes, 0);
    
    // Clear form
    document.getElementById('exercise-minutes').value = '';
    document.getElementById('exercise-type').value = '';
    
    // Show success message
    const typeEmoji = getExerciseEmoji(type);
    showToast(`${typeEmoji} ${minutes} min de ${type} registrados! Total hoje: ${totalMinutes} min`, 'success');
    
    if (totalMinutes >= 30) {
        showToast('🎉 Parabéns! Você atingiu sua meta diária de exercícios!', 'success');
    }
}

function getExerciseEmoji(type) {
    const emojis = {
        'caminhada': '🚶',
        'corrida': '🏃',
        'yoga': '🧘',
        'academia': '💪',
        'natacao': '🏊',
        'ciclismo': '🚴'
    };
    return emojis[type] || '🏃';
}

// Meditation Timer Functions
let meditationInterval = null;
let meditationTimeLeft = 300; // 5 minutes in seconds

function setMeditationTime(minutes) {
    meditationTimeLeft = minutes * 60;
    localStorage.setItem('meditation_time', minutes.toString());
    updateMeditationDisplay();
    
    // Update button states
    document.querySelectorAll('.timer-controls .btn').forEach(btn => {
        btn.classList.remove('active');
    });
    document.querySelector(`[onclick="setMeditationTime(${minutes})"]`).classList.add('active');
}

function toggleMeditation() {
    const btn = document.getElementById('meditation-btn');
    
    if (meditationInterval) {
        // Stop meditation
        clearInterval(meditationInterval);
        meditationInterval = null;
        btn.innerHTML = '<i class="fas fa-play"></i> Iniciar';
        btn.classList.remove('btn-danger');
        btn.classList.add('btn-primary');
    } else {
        // Start meditation
        meditationInterval = setInterval(() => {
            meditationTimeLeft--;
            updateMeditationDisplay();
            
            if (meditationTimeLeft <= 0) {
                // Meditation finished
                clearInterval(meditationInterval);
                meditationInterval = null;
                btn.innerHTML = '<i class="fas fa-play"></i> Iniciar';
                btn.classList.remove('btn-danger');
                btn.classList.add('btn-primary');
                
                // Reset timer
                const savedTime = parseInt(localStorage.getItem('meditation_time') || '5');
                setMeditationTime(savedTime);
                
                showToast('🧘 Sessão de meditação concluída! Parabéns!', 'success');
                
                // Save meditation session
                const today = getTodayDate();
                const sessions = JSON.parse(localStorage.getItem(`meditation_${today}`) || '[]');
                sessions.push({
                    duration: savedTime,
                    timestamp: new Date().toISOString()
                });
                localStorage.setItem(`meditation_${today}`, JSON.stringify(sessions));
            }
        }, 1000);
        
        btn.innerHTML = '<i class="fas fa-stop"></i> Parar';
        btn.classList.remove('btn-primary');
        btn.classList.add('btn-danger');
    }
}

function updateMeditationDisplay() {
    const minutes = Math.floor(meditationTimeLeft / 60);
    const seconds = meditationTimeLeft % 60;
    const display = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    
    const timerElement = document.getElementById('meditation-timer');
    if (timerElement) {
        timerElement.textContent = display;
    }
}

// User Filter Functions
function filterUsers(type) {
    // Update active filter button
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    // Filter users (this would integrate with the existing loadUsers function)
    loadUsers('', type);
}

// Health Post Functions
function showCreateHealthPost() {
    showCreatePost(); // Use existing function but with health-specific form
}

// Health Resources Functions
function showHealthTips() {
    showModal('Dicas de Saúde', 
        'Aqui estão algumas dicas importantes:\n\n' +
        '• Beba pelo menos 8 copos de água por dia\n' +
        '• Pratique 30 minutos de exercício diariamente\n' +
        '• Durma 7-8 horas por noite\n' +
        '• Mantenha uma alimentação equilibrada\n' +
        '• Pratique técnicas de relaxamento\n' +
        '• Faça check-ups médicos regulares'
    );
}

function showFirstAid() {
    showModal('Primeiros Socorros', 
        'Informações básicas de primeiros socorros:\n\n' +
        '• Em caso de emergência, ligue 192 (SAMU)\n' +
        '• Para parada cardíaca: RCP (30 compressões + 2 ventilações)\n' +
        '• Para engasgo: Manobra de Heimlich\n' +
        '• Para queimaduras: água fria por 10-15 minutos\n' +
        '• Para cortes: pressão direta para estancar sangramento\n\n' +
        'IMPORTANTE: Procure sempre ajuda médica profissional!'
    );
}

function showHealthyRecipes() {
    showModal('Receitas Saudáveis', 
        'Algumas ideias de receitas saudáveis:\n\n' +
        '🥗 Salada Colorida:\n' +
        '- Folhas verdes, tomate, cenoura, pepino\n' +
        '- Azeite, limão e ervas\n\n' +
        '🍲 Sopa de Legumes:\n' +
        '- Abóbora, cenoura, abobrinha, cebola\n' +
        '- Temperos naturais\n\n' +
        '🐟 Peixe Grelhado:\n' +
        '- Salmão ou tilápia\n' +
        '- Legumes no vapor\n\n' +
        '🥤 Smoothie Verde:\n' +
        '- Espinafre, banana, maçã, água de coco'
    );
}

function showMentalHealth() {
    showModal('Saúde Mental', 
        'Cuidando da sua saúde mental:\n\n' +
        '🧠 Sinais de alerta:\n' +
        '- Tristeza persistente\n' +
        '- Ansiedade excessiva\n' +
        '- Mudanças no sono/apetite\n' +
        '- Isolamento social\n\n' +
        '💚 Como cuidar:\n' +
        '- Pratique mindfulness\n' +
        '- Mantenha conexões sociais\n' +
        '- Exercite-se regularmente\n' +
        '- Busque ajuda profissional quando necessário\n\n' +
        '📞 CVV: 188 (24h, gratuito)'
    );
}

// Emergency Functions
function showEmergencyContacts() {
    showModal('Contatos de Emergência', 
        'Números importantes para emergências:\n\n' +
        '🚑 SAMU: 192\n' +
        '🚒 Bombeiros: 193\n' +
        '👮 Polícia: 190\n' +
        '💚 CVV (Apoio emocional): 188\n' +
        '🏥 Disque Saúde: 136\n' +
        '☠️ Centro de Intoxicações: 0800 722 6001\n\n' +
        'Mantenha esses números sempre acessíveis!'
    );
}

// Enhanced User Loading with Health Filters
async function loadUsersWithHealthFilter(searchTerm = '', userType = 'all') {
    try {
        showLoader();
        
        let endpoint = '/usuarios';
        if (searchTerm) {
            endpoint += `?search=${encodeURIComponent(searchTerm)}`;
        }
        
        const response = await apiRequest(endpoint);
        let users = response;
        
        // Filter by user type if specified
        if (userType !== 'all') {
            users = users.filter(user => {
                // This would require adding user type to the user model
                return user.tipo === userType;
            });
        }
        
        displayUsersWithHealthInfo(users);
        
    } catch (error) {
        console.error('Erro ao carregar usuários:', error);
        showToast('Erro ao carregar usuários', 'error');
    } finally {
        hideLoader();
    }
}

function displayUsersWithHealthInfo(users) {
    const usersList = document.getElementById('users-list');
    
    if (!users || users.length === 0) {
        usersList.innerHTML = '<p class="no-users">Nenhum usuário encontrado.</p>';
        return;
    }
    
    usersList.innerHTML = users.map(user => {
        const isFollowing = followingUsersCache.has(user.ID);
        const userTypeIcon = getUserTypeIcon(user.tipo || 'paciente');
        const userTypeBadge = getUserTypeBadge(user.tipo || 'paciente');
        
        return `
            <div class="user-card health-user-card" data-user-id="${user.ID}">
                <div class="user-avatar">
                    ${user.FotoPerfil ? 
                        `<img src="${API_BASE_URL}/uploads/perfil/${user.FotoPerfil}" alt="Foto de ${user.Nome}" onerror="this.style.display='none'; this.nextElementSibling.style.display='block';">
                         <i class="fas fa-user-circle" style="display: none;"></i>` :
                        `<i class="fas fa-user-circle"></i>`
                    }
                    <div class="user-type-indicator">${userTypeIcon}</div>
                </div>
                <div class="user-info">
                    <h3>${user.Nome}</h3>
                    <p>@${user.Nick}</p>
                    <div class="user-badges">
                        ${userTypeBadge}
                        ${user.verificado ? '<span class="badge badge-verified">✓ Verificado</span>' : ''}
                    </div>
                </div>
                <div class="user-actions">
                    ${user.ID !== getCurrentUserId() ? `
                        <button class="btn ${isFollowing ? 'btn-secondary' : 'btn-primary'} user-action-btn" 
                                data-user-id="${user.ID}" 
                                data-action="${isFollowing ? 'unfollow' : 'follow'}">
                            <i class="fas ${isFollowing ? 'fa-user-minus' : 'fa-user-plus'}"></i>
                            ${isFollowing ? 'Deixar de seguir' : 'Seguir'}
                        </button>
                    ` : '<span class="current-user-badge">Você</span>'}
                </div>
            </div>
        `;
    }).join('');
    
    addUserActionListeners();
}

function getUserTypeIcon(type) {
    const icons = {
        'profissional': '👨‍⚕️',
        'paciente': '👤',
        'cuidador': '🤝',
        'estudante': '📚'
    };
    return icons[type] || '👤';
}

function getUserTypeBadge(type) {
    const badges = {
        'profissional': '<span class="badge badge-professional">👨‍⚕️ Profissional</span>',
        'paciente': '<span class="badge badge-user">👤 Paciente</span>',
        'cuidador': '<span class="badge badge-user">🤝 Cuidador</span>',
        'estudante': '<span class="badge badge-user">📚 Estudante</span>'
    };
    return badges[type] || '<span class="badge badge-user">👤 Usuário</span>';
}

// Enhanced Dashboard Stats for Health Theme
async function loadHealthDashboardStats() {
    try {
        // Load existing stats
        await loadDashboardStats();
        
        // Add health-specific stats
        const professionalsCount = 0; // This would come from API
        const postsCount = 0; // This would come from API
        const interactionsCount = 0; // This would come from API
        
        // Update health stats
        const professionalsElement = document.getElementById('professionals-count');
        const postsElement = document.getElementById('posts-count');
        const interactionsElement = document.getElementById('interactions-count');
        
        if (professionalsElement) professionalsElement.textContent = professionalsCount;
        if (postsElement) postsElement.textContent = postsCount;
        if (interactionsElement) interactionsElement.textContent = interactionsCount;
        
    } catch (error) {
        console.error('Erro ao carregar estatísticas de saúde:', error);
    }
}

// Initialize health theme when page loads
document.addEventListener('DOMContentLoaded', function() {
    // Override the showHome function to load health stats
    const originalShowHome = window.showHome;
    window.showHome = function() {
        originalShowHome();
        loadHealthDashboardStats();
    };
    
    // Add wellness navigation if user is logged in
    if (getAuthToken()) {
        // The wellness button is already in the HTML
    }
});