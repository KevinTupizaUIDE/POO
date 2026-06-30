// State management
let currentFilter = 'all'; // 'all', 'tech', 'cheap', 'avail'
let booksData = [];

// DOM Elements
const booksGrid = document.getElementById('books-grid');
const searchInput = document.getElementById('search-input');
const searchBtn = document.getElementById('search-btn');
const statTotal = document.getElementById('stat-total');
const pageTitle = document.getElementById('page-title');
const pageSubtitle = document.getElementById('page-subtitle');
const activeFiltersContainer = document.getElementById('active-filters-container');

// Sidebar menu buttons
const btnAllMenu = document.getElementById('btn-all-menu');
const btnTechMenu = document.getElementById('btn-tech-menu');
const btnCheapMenu = document.getElementById('btn-cheap-menu');
const btnAvailMenu = document.getElementById('btn-avail-menu');

// Modal Elements
const detailsModal = document.getElementById('details-modal');
const closeModalBtn = document.getElementById('close-modal');
const modalTitle = document.getElementById('modal-title');
const modalBody = document.getElementById('modal-body');

// Initialize event listeners
document.addEventListener('DOMContentLoaded', () => {
    fetchBooks();
    setupNavigation();
    setupSearch();
    setupModal();
});

// Setup sidebar navigation links
function setupNavigation() {
    const navItems = [
        { btn: btnAllMenu, filter: 'all', title: 'Catálogo Completo', subtitle: 'Explora todos los libros electrónicos y filtra usando programación funcional.' },
        { btn: btnTechMenu, filter: 'tech', title: 'Libros de Tecnología', subtitle: 'Filtro funcional aplicado para el género de Tecnología.' },
        { btn: btnCheapMenu, filter: 'cheap', title: 'Libros Económicos', subtitle: 'Libros electrónicos con un precio menor o igual a $30.00.' },
        { btn: btnAvailMenu, filter: 'avail', title: 'Disponibles para Descarga', subtitle: 'Mostrando libros que tienen disponibilidad de stock activa.' }
    ];

    navItems.forEach(item => {
        item.btn.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Remove active class from all navigation items
            document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
            
            // Add active class to clicked item
            item.btn.classList.add('active');
            
            // Update title and subtitles
            pageTitle.innerText = item.title;
            pageSubtitle.innerText = item.subtitle;
            
            currentFilter = item.filter;
            
            // Reset search input
            searchInput.value = '';
            
            // Fetch and apply current filter
            fetchBooks();
        });
    });
}

// Fetch catalog from API
async function fetchBooks() {
    showLoading();
    let url = '/api/catalogo';
    
    // Append query parameters based on active filter
    const params = new URLSearchParams();
    if (currentFilter === 'tech') {
        params.append('genre', 'Tecnología');
    } else if (currentFilter === 'cheap') {
        params.append('max_price', '30.00');
    } else if (currentFilter === 'avail') {
        params.append('available', 'true');
    }

    if (params.toString()) {
        url += `?${params.toString()}`;
    }

    try {
        const response = await fetch(url);
        if (!response.ok) throw new Error('Error al cargar catálogo');
        
        booksData = await response.json();
        
        // Handle empty response
        if (!booksData) booksData = [];
        
        renderBooks(booksData);
        updateStats(booksData.length);
        renderActiveFilters();
    } catch (error) {
        console.error(error);
        renderError('No se pudo conectar al servidor o cargar el catálogo de libros.');
    }
}

// Render active filters badges
function renderActiveFilters() {
    activeFiltersContainer.innerHTML = '';
    
    if (currentFilter !== 'all') {
        let label = '';
        if (currentFilter === 'tech') label = 'Género: Tecnología';
        else if (currentFilter === 'cheap') label = 'Precio ≤ $30.00';
        else if (currentFilter === 'avail') label = 'Disponible';

        const badge = document.createElement('div');
        badge.className = 'filter-badge';
        badge.innerHTML = `
            <span>${label}</span>
            <span class="clear-badge" id="clear-filter-x">&times;</span>
        `;
        
        activeFiltersContainer.appendChild(badge);
        
        document.getElementById('clear-filter-x').addEventListener('click', () => {
            btnAllMenu.click();
        });
    }
}

// Render list of book cards
function renderBooks(books) {
    booksGrid.innerHTML = '';
    
    if (books.length === 0) {
        booksGrid.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <h3>No se encontraron libros</h3>
                <p>No hay registros que coincidan con el criterio seleccionado actualmente.</p>
            </div>
        `;
        return;
    }

    books.forEach((book, index) => {
        const card = document.createElement('div');
        card.className = 'book-card';
        card.style.animationDelay = `${index * 0.05}s`;
        
        const availabilityText = book.IsAvailable ? 'Disponible para descarga inmediata' : 'No disponible temporalmente';
        const availabilityClass = book.IsAvailable ? 'available' : 'unavailable';
        
        card.innerHTML = `
            <div class="book-header">
                <span class="genre-tag">${book.Genre}</span>
                <span class="book-id">#ID ${book.ID}</span>
            </div>
            <div>
                <h3 class="book-title">${book.Title}</h3>
                <p class="book-author">Por ${book.Author}</p>
            </div>
            <div class="book-footer">
                <div class="price-container">
                    <span class="price-label">Precio</span>
                    <span class="price-value">$${book.Price.toFixed(2)}</span>
                </div>
                <div class="availability-dot ${availabilityClass}" data-tooltip="${availabilityText}"></div>
            </div>
        `;
        
        // Open details modal on card click
        card.addEventListener('click', () => {
            openBookModal(book);
        });

        booksGrid.appendChild(card);
    });
}

// Setup search bar
function setupSearch() {
    const performSearch = async () => {
        const idQuery = searchInput.value.trim();
        if (!idQuery) {
            fetchBooks();
            return;
        }

        // Check if query is numeric
        if (isNaN(idQuery)) {
            renderError('Por favor, ingresa un número de ID válido para buscar.');
            return;
        }

        showLoading();
        try {
            const response = await fetch(`/api/catalogo/buscar?id=${idQuery}`);
            if (response.status === 404) {
                booksGrid.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">🔍</div>
                        <h3>Libro no encontrado</h3>
                        <p>No se encontró ningún libro con el ID <strong>${idQuery}</strong> en la base de datos.</p>
                        <button class="primary-btn" style="margin-top: 1rem;" onclick="document.getElementById('search-input').value=''; fetchBooks();">Mostrar Todo</button>
                    </div>
                `;
                updateStats(0);
                return;
            }
            if (!response.ok) throw new Error('Error al buscar por ID');

            const book = await response.json();
            renderBooks([book]);
            updateStats(1);
            
            // Set active navigation to "All" since we searched
            document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
            btnAllMenu.classList.add('active');
            pageTitle.innerText = `Búsqueda por ID: ${idQuery}`;
            pageSubtitle.innerText = `Resultado obtenido directamente desde el controlador CRUD de GORM.`;

            // Active filters indicator
            activeFiltersContainer.innerHTML = `
                <div class="filter-badge">
                    <span>Búsqueda: ID ${idQuery}</span>
                    <span class="clear-badge" onclick="document.getElementById('search-input').value=''; document.getElementById('btn-all-menu').click();">&times;</span>
                </div>
            `;
        } catch (error) {
            console.error(error);
            renderError('Ocurrió un error al realizar la búsqueda en la API.');
        }
    };

    searchBtn.addEventListener('click', performSearch);
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            performSearch();
        }
    });
}

// Open modal showing details and encapsulation demonstration
function openBookModal(book) {
    modalTitle.innerText = book.Title;
    
    // Simulate encapsulation value locally (if not persisted, we show it client-side)
    if (!book.hasOwnProperty('estadoInternoLocal')) {
        book.estadoInternoLocal = "Sin auditar"; // Default from encapsulation requirement
    }

    const availableText = book.IsAvailable ? 'Sí (Activo)' : 'No (Inactivo)';
    const availableColor = book.IsAvailable ? 'var(--text-success)' : 'var(--text-danger)';

    modalBody.innerHTML = `
        <div class="detail-row">
            <span class="detail-label">ID Registro GORM</span>
            <span class="detail-value highlight">#${book.ID}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Autor del Libro</span>
            <span class="detail-value">${book.Author}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Género / Categoría</span>
            <span class="detail-value">${book.Genre}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Precio Unitario</span>
            <span class="detail-value">$${book.Price.toFixed(2)}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Disponible Descarga</span>
            <span class="detail-value" style="color: ${availableColor}">${availableText}</span>
        </div>
        <div class="detail-row" style="background: rgba(0, 210, 255, 0.02); padding: 0.85rem; border-radius: 8px; border: 1px solid rgba(0, 210, 255, 0.1);">
            <span class="detail-label">Estado Interno (Getter)</span>
            <span class="detail-value" id="modal-estado-interno" style="color: var(--accent-blue); font-weight: bold;">${book.estadoInternoLocal}</span>
        </div>
        
        <div class="interactive-box">
            <label>Demostración de Encapsulación (Setter):</label>
            <div class="interactive-input-group">
                <input type="text" id="mutation-input" placeholder="Ingresa un nuevo estado...">
                <button class="primary-btn" id="mutate-btn">Modificar</button>
            </div>
        </div>
    `;

    detailsModal.classList.add('active');

    // Add encapsulation setter handler
    const mutateBtn = document.getElementById('mutate-btn');
    const mutationInput = document.getElementById('mutation-input');
    const estadoDisplay = document.getElementById('modal-estado-interno');

    mutateBtn.addEventListener('click', () => {
        const newValue = mutationInput.value.trim();
        if (newValue) {
            // Apply encapsulation SetEstadoInterno logic
            book.estadoInternoLocal = newValue;
            estadoDisplay.innerText = newValue;
            mutationInput.value = '';
            
            // Add subtle pulse effect
            estadoDisplay.style.transition = 'none';
            estadoDisplay.style.transform = 'scale(1.15)';
            estadoDisplay.style.color = '#ffffff';
            setTimeout(() => {
                estadoDisplay.style.transition = 'all 0.3s ease';
                estadoDisplay.style.transform = 'scale(1)';
                estadoDisplay.style.color = 'var(--accent-blue)';
            }, 150);
        }
    });
}

// Modal handling
function setupModal() {
    closeModalBtn.addEventListener('click', () => {
        detailsModal.classList.remove('active');
    });

    detailsModal.addEventListener('click', (e) => {
        if (e.target === detailsModal) {
            detailsModal.classList.remove('active');
        }
    });
}

// UI State Helpers
function showLoading() {
    booksGrid.innerHTML = `
        <div class="loading-spinner">
            <div class="spinner"></div>
            <p>Consultando base de datos y aplicando filtros funcionales...</p>
        </div>
    `;
}

function renderError(message) {
    booksGrid.innerHTML = `
        <div class="empty-state">
            <div class="empty-state-icon" style="color: var(--text-danger);">⚠️</div>
            <h3>Error de Conexión</h3>
            <p>${message}</p>
            <button class="primary-btn" style="margin-top: 1rem; background-color: var(--text-danger); color: #ffffff;" onclick="fetchBooks()">Reintentar</button>
        </div>
    `;
    updateStats(0);
}

function updateStats(count) {
    statTotal.innerText = count;
}
