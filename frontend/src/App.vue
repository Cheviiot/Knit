<template>
  <div class="app">
    <!-- Titlebar -->
    <header class="titlebar" style="--wails-draggable:drag">
      <div class="titlebar-left">
        <div class="brand">
          <img src="./assets/Knit.png" alt="Knit" class="brand-icon" />
          <span class="brand-name">Knit</span>
        </div>
      </div>
      
      <div class="titlebar-center" style="--wails-draggable:no-drag">
        <div class="search-box">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="M21 21l-4.35-4.35"/>
          </svg>
          <input 
            v-model="searchQuery" 
            type="text" 
            class="search-input" 
            :placeholder="contentType === 'movies' ? 'Поиск фильмов...' : 'Поиск сериалов...'"
            @keyup.enter="searchContent"
          />
          <button v-if="searchQuery" class="search-clear" @click="clearSearch">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
      
      <div class="titlebar-right">
        <div class="content-type-toggle" style="--wails-draggable:no-drag">
          <button 
            class="toggle-btn" 
            :class="{ active: contentType === 'movies' }" 
            @click="switchContentType('movies')"
            title="Фильмы"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18"/>
              <line x1="7" y1="2" x2="7" y2="22"/>
              <line x1="17" y1="2" x2="17" y2="22"/>
              <line x1="2" y1="12" x2="22" y2="12"/>
            </svg>
          </button>
          <button 
            class="toggle-btn" 
            :class="{ active: contentType === 'tv' }" 
            @click="switchContentType('tv')"
            title="Сериалы"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="15" rx="2" ry="2"/>
              <polyline points="17 2 12 7 7 2"/>
            </svg>
          </button>
        </div>
        
        <button class="titlebar-btn" @click="showSettings = true" title="Настройки" style="--wails-draggable:no-drag">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
          </svg>
        </button>
        
        <div class="titlebar-divider"></div>
        
        <div class="window-controls" style="--wails-draggable:no-drag">
          <button class="window-btn" @click="minimizeWindow">
            <svg width="10" height="10" viewBox="0 0 12 12"><rect y="5" width="12" height="1" fill="currentColor"/></svg>
          </button>
          <button class="window-btn" @click="toggleMaximize">
            <svg width="10" height="10" viewBox="0 0 12 12"><rect x="1" y="1" width="10" height="10" fill="none" stroke="currentColor" stroke-width="1"/></svg>
          </button>
          <button class="window-btn window-close" @click="closeWindow">
            <svg width="10" height="10" viewBox="0 0 12 12"><path d="M1 1l10 10M11 1L1 11" stroke="currentColor" stroke-width="1.2"/></svg>
          </button>
        </div>
      </div>
    </header>
    
    <!-- Main Content -->
    <main class="main">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка...</p>
      </div>
      
      <!-- Movies Grid -->
      <div v-else-if="contentType === 'movies' && movies.length > 0" class="movies-section">
        <h2 class="section-title">
          {{ isSearchResult ? 'Результаты поиска' : 'Популярные фильмы' }}
          <span class="count">({{ movies.length }})</span>
        </h2>
        
        <div class="movies-grid">
          <div 
            v-for="movie in movies" 
            :key="movie.id" 
            class="movie-card group"
            @click="selectMovie(movie)"
          >
            <div class="movie-poster">
              <img 
                v-if="movie.poster_path" 
                :src="getImageUrl(movie.poster_path, 'w500')" 
                :alt="movie.title"
                loading="lazy"
                class="transition-transform duration-300 group-hover:scale-105"
              />
              <div v-else class="no-poster">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18"/>
                  <line x1="7" y1="2" x2="7" y2="22"/>
                  <line x1="17" y1="2" x2="17" y2="22"/>
                  <line x1="2" y1="12" x2="22" y2="12"/>
                  <line x1="2" y1="7" x2="7" y2="7"/>
                  <line x1="2" y1="17" x2="7" y2="17"/>
                  <line x1="17" y1="17" x2="22" y2="17"/>
                  <line x1="17" y1="7" x2="22" y2="7"/>
                </svg>
              </div>
              <!-- Hover overlay with gradient -->
              <div class="absolute inset-0 bg-linear-to-t from-black/80 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <!-- Rating badge -->
              <div class="movie-rating" v-if="movie.vote_average">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                </svg>
                {{ movie.vote_average.toFixed(1) }}
              </div>
            </div>
            <div class="movie-info">
              <h3 class="movie-title">{{ movie.title }}</h3>
              <p class="movie-year">{{ getYear(movie.release_date) }}</p>
            </div>
          </div>
        </div>
        
        <!-- Loading More Indicator -->
        <div v-if="loadingMore" class="loading-more">
          <div class="spinner"></div>
          <p>Загрузка...</p>
        </div>
      </div>
      
      <!-- TV Shows Grid -->
      <div v-else-if="contentType === 'tv' && tvShows.length > 0" class="movies-section">
        <h2 class="section-title">
          {{ isSearchResult ? 'Результаты поиска' : 'Популярные сериалы' }}
          <span class="count">({{ tvShows.length }})</span>
        </h2>
        
        <div class="movies-grid">
          <div 
            v-for="show in tvShows" 
            :key="show.id" 
            class="movie-card group"
            @click="selectTVShow(show)"
          >
            <div class="movie-poster">
              <img 
                v-if="show.poster_path" 
                :src="getImageUrl(show.poster_path, 'w500')" 
                :alt="show.name"
                loading="lazy"
                class="transition-transform duration-300 group-hover:scale-105"
              />
              <div v-else class="no-poster">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <rect x="2" y="7" width="20" height="15" rx="2" ry="2"/>
                  <polyline points="17 2 12 7 7 2"/>
                </svg>
              </div>
              <!-- Hover overlay with gradient -->
              <div class="absolute inset-0 bg-linear-to-t from-black/80 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <!-- Rating badge -->
              <div class="movie-rating" v-if="show.vote_average">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                </svg>
                {{ show.vote_average.toFixed(1) }}
              </div>
            </div>
            <div class="movie-info">
              <h3 class="movie-title">{{ show.name }}</h3>
              <p class="movie-year">{{ getYear(show.first_air_date) }}</p>
            </div>
          </div>
        </div>
        
        <!-- Loading More Indicator -->
        <div v-if="loadingMore" class="loading-more">
          <div class="spinner"></div>
          <p>Загрузка...</p>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-else class="empty-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
        <h3>Начните поиск</h3>
        <p>{{ contentType === 'movies' ? 'Введите название фильма в строке поиска' : 'Введите название сериала в строке поиска' }}</p>
      </div>
    </main>
    
    <!-- Movie Details Modal -->
    <Transition name="modal">
      <div v-if="selectedMovie" class="modal-overlay" @click.self="closeModal">
        <div class="movie-modal">
          <!-- Backdrop на весь фон -->
          <div class="modal-backdrop">
            <img v-if="selectedMovie.backdrop_path" :src="getImageUrl(selectedMovie.backdrop_path, 'w1280')" />
            <div class="backdrop-overlay"></div>
          </div>
          
          <button class="modal-close" @click="closeModal">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
          
          <div class="modal-content">
            <div class="modal-poster">
              <img 
                v-if="selectedMovie.poster_path" 
                :src="getImageUrl(selectedMovie.poster_path, 'w500')"
              />
            </div>
            
            <div class="modal-info">
              <h2 class="modal-title">{{ selectedMovie.title }}</h2>
              <p class="modal-original" v-if="selectedMovie.original_title !== selectedMovie.title">
                {{ selectedMovie.original_title }}
              </p>
              
              <div class="modal-meta">
                <span class="meta-item" v-if="selectedMovie.release_date">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                    <line x1="16" y1="2" x2="16" y2="6"/>
                    <line x1="8" y1="2" x2="8" y2="6"/>
                    <line x1="3" y1="10" x2="21" y2="10"/>
                  </svg>
                  {{ formatDate(selectedMovie.release_date) }}
                </span>
                <span class="meta-item rating" v-if="selectedMovie.vote_average">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                  </svg>
                  {{ selectedMovie.vote_average.toFixed(1) }}
                </span>
              </div>
              
              <p class="modal-overview">{{ selectedMovie.overview || 'Описание отсутствует' }}</p>
            </div>
          </div>
          
          <!-- Torrents Section -->
          <div v-if="torrentsLoading || torrents.length > 0 || torrentsError" class="torrents-section">
            <h3 class="torrents-title">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Торренты
              <span v-if="torrents.length" class="count">({{ torrents.length }})</span>
            </h3>
            
            <div v-if="torrentsLoading" class="torrents-loading">
              <div class="spinner"></div>
              <span>Поиск торрентов...</span>
            </div>
            
            <div v-else-if="torrentsError" class="torrents-error">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ torrentsError }}
            </div>
            
            <div v-else class="torrents-list">
              <div v-for="(torrent, index) in torrents" :key="index" class="torrent-item">
                <div class="torrent-info">
                  <h4 class="torrent-title">{{ torrent.title }}</h4>
                  <div class="torrent-meta">
                    <span class="torrent-size">{{ torrent.sizeStr }}</span>
                    <span class="torrent-tracker">{{ torrent.tracker }}</span>
                  </div>
                </div>
                <div class="torrent-stats">
                  <span class="seeders">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="18 15 12 9 6 15"/>
                    </svg>
                    {{ torrent.seeders }}
                  </span>
                  <span class="leechers">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="6 9 12 15 18 9"/>
                    </svg>
                    {{ torrent.leechers }}
                  </span>
                </div>
                <div class="torrent-actions">
                  <button v-if="torrent.magnetUri || torrent.link" class="btn btn-accent btn-sm" @click="downloadTorrent(torrent)">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                      <polyline points="7 10 12 15 17 10"/>
                      <line x1="12" y1="15" x2="12" y2="3"/>
                    </svg>
                    В FDM
                  </button>
                  <button v-if="torrent.magnetUri" class="btn btn-secondary btn-sm" @click="copyMagnet(torrent.magnetUri)">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/>
                      <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/>
                    </svg>
                    Magnet
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- TV Show Details Modal -->
    <Transition name="modal">
      <div v-if="selectedTVShow" class="modal-overlay" @click.self="closeTVModal">
        <div class="movie-modal">
          <!-- Backdrop на весь фон -->
          <div class="modal-backdrop">
            <img v-if="selectedTVShow.backdrop_path" :src="getImageUrl(selectedTVShow.backdrop_path, 'w1280')" />
            <div class="backdrop-overlay"></div>
          </div>
          
          <button class="modal-close" @click="closeTVModal">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
          
          <div class="modal-content">
            <div class="modal-poster">
              <img 
                v-if="selectedTVShow.poster_path" 
                :src="getImageUrl(selectedTVShow.poster_path, 'w500')"
              />
            </div>
            
            <div class="modal-info">
              <h2 class="modal-title">{{ selectedTVShow.name }}</h2>
              <p class="modal-original" v-if="selectedTVShow.original_name !== selectedTVShow.name">
                {{ selectedTVShow.original_name }}
              </p>
              
              <div class="modal-meta">
                <span class="meta-item" v-if="selectedTVShow.first_air_date">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                    <line x1="16" y1="2" x2="16" y2="6"/>
                    <line x1="8" y1="2" x2="8" y2="6"/>
                    <line x1="3" y1="10" x2="21" y2="10"/>
                  </svg>
                  {{ formatDate(selectedTVShow.first_air_date) }}
                </span>
                <span class="meta-item rating" v-if="selectedTVShow.vote_average">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                  </svg>
                  {{ selectedTVShow.vote_average.toFixed(1) }}
                </span>
              </div>
              
              <p class="modal-overview">{{ selectedTVShow.overview || 'Описание отсутствует' }}</p>
            </div>
          </div>
          
          <!-- Torrents Section -->
          <div v-if="torrentsLoading || torrents.length > 0 || torrentsError" class="torrents-section">
            <h3 class="torrents-title">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Торренты
              <span v-if="torrents.length" class="count">({{ torrents.length }})</span>
            </h3>
            
            <div v-if="torrentsLoading" class="torrents-loading">
              <div class="spinner"></div>
              <span>Поиск торрентов...</span>
            </div>
            
            <div v-else-if="torrentsError" class="torrents-error">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ torrentsError }}
            </div>
            
            <div v-else class="torrents-list">
              <div v-for="(torrent, index) in torrents" :key="index" class="torrent-item">
                <div class="torrent-info">
                  <h4 class="torrent-title">{{ torrent.title }}</h4>
                  <div class="torrent-meta">
                    <span class="torrent-size">{{ torrent.sizeStr }}</span>
                    <span class="torrent-seeds">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 2L2 12l10 10 10-10L12 2z"/>
                      </svg>
                      {{ torrent.seeders }}
                    </span>
                    <span class="torrent-tracker" v-if="torrent.tracker">{{ torrent.tracker }}</span>
                  </div>
                </div>
                <div class="torrent-actions">
                  <button class="torrent-btn btn-copy" @click="copyMagnet(torrent.magnetUri)" title="Копировать magnet">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                      <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                    </svg>
                  </button>
                  <button class="torrent-btn btn-download" @click="downloadTorrent(torrent)" title="Скачать через FDM">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                      <polyline points="7 10 12 15 17 10"/>
                      <line x1="12" y1="15" x2="12" y2="3"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- Downloads Drawer -->
    <Transition name="drawer">
      <div v-if="showDownloads" class="downloads-overlay" @click.self="showDownloads = false">
        <div class="downloads-drawer">
          <div class="drawer-header">
            <div class="drawer-title">
              <div class="drawer-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                  <polyline points="7 10 12 15 17 10"/>
                  <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
              </div>
              <div>
                <h3>Загрузки</h3>
                <span class="drawer-subtitle">Free Download Manager</span>
              </div>
            </div>
            <div class="drawer-actions">
              <button class="drawer-btn" @click="openDownloadsFolder" title="Открыть папку">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/>
                </svg>
              </button>
              <button class="drawer-btn drawer-close" @click="showDownloads = false">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
          
          <div class="drawer-content">
            <div v-if="downloads.length === 0" class="drawer-empty">
              <div class="empty-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                  <polyline points="7 10 12 15 17 10"/>
                  <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
              </div>
              <p>Список пуст</p>
              <span>Торренты будут отправлены в FDM</span>
            </div>
            
            <div v-else class="drawer-list">
              <div v-for="dl in downloads" :key="dl.id" class="dl-card">
                <div class="dl-status-indicator" :class="'indicator-' + dl.status.replace('sent_to_', '')"></div>
                <div class="dl-content">
                  <div class="dl-name">{{ dl.name }}</div>
                  <div class="dl-info">
                    <span class="dl-status" :class="'status-' + dl.status.replace('sent_to_', '')">
                      {{ getStatusText(dl.status) }}
                    </span>
                    <span v-if="dl.client" class="dl-client">• {{ dl.client }}</span>
                  </div>
                </div>
                <div class="dl-actions">
                  <button v-if="dl.status.startsWith('sent_to_')" class="dl-btn" @click="openFDM" title="Открыть">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6"/>
                      <polyline points="15 3 21 3 21 9"/>
                      <line x1="10" y1="14" x2="21" y2="3"/>
                    </svg>
                  </button>
                  <button class="dl-btn dl-btn-remove" @click="removeDownload(dl.id)" title="Удалить">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- Downloads FAB -->
    <button v-if="downloads.length > 0" class="downloads-fab" :class="{ 'fab-active': showDownloads }" @click="showDownloads = !showDownloads">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
      <span v-if="downloads.length > 0" class="fab-badge">{{ downloads.length }}</span>
    </button>
    
    <!-- Settings Modal -->
    <Transition name="modal">
      <div v-if="showSettings" class="modal-overlay" @click.self="showSettings = false">
        <div class="settings-modal">
          <div class="settings-header">
            <h2>Настройки</h2>
            <button class="modal-close" @click="showSettings = false">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <div class="settings-content">
            <div class="settings-group">
              <label class="settings-label">Парсер (Jackett сервер)</label>
              <div class="settings-row">
                <div class="dropdown dropdown-full">
                  <button class="dropdown-trigger dropdown-input" @click="toggleDropdown('settingsServer')">
                    <span class="dropdown-value">{{ settingsServerName }}</span>
                    <svg class="dropdown-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="6 9 12 15 18 9"/>
                    </svg>
                  </button>
                  <Transition name="dropdown">
                    <div v-if="openDropdown === 'settingsServer'" class="dropdown-menu">
                      <button 
                        v-for="server in servers" 
                        :key="server.id" 
                        class="dropdown-item"
                        :class="{ 'active': settings.selectedServer === server.id }"
                        @click="selectSettingsServer(server.id)"
                      >
                        {{ server.name }}
                        <svg v-if="settings.selectedServer === server.id" class="dropdown-check" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </button>
                    </div>
                  </Transition>
                </div>
                <button 
                  class="btn btn-check" 
                  :class="{ 'checking': checkingServer }"
                  @click="checkServer"
                  :disabled="checkingServer"
                >
                  <span v-if="checkingServer" class="spinner-small"></span>
                  <span v-else>Проверить</span>
                </button>
              </div>
              <div class="settings-status" v-if="serverStatus !== null">
                <span :class="serverStatus ? 'status-online' : 'status-offline'">
                  {{ serverStatus ? '✓ Сервер доступен' : '✗ Сервер недоступен' }}
                </span>
              </div>
              <span class="settings-hint">Сервер для поиска торрентов</span>
            </div>
            
            <div class="settings-group">
              <label class="settings-label">TMDB прокси</label>
              <div class="settings-row">
                <div class="dropdown dropdown-full">
                  <button class="dropdown-trigger dropdown-input" @click="toggleDropdown('settingsProxy')">
                    <span class="dropdown-value">{{ settingsProxyName }}</span>
                    <svg class="dropdown-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="6 9 12 15 18 9"/>
                    </svg>
                  </button>
                  <Transition name="dropdown">
                    <div v-if="openDropdown === 'settingsProxy'" class="dropdown-menu">
                      <button 
                        v-for="proxy in tmdbProxies" 
                        :key="proxy.id" 
                        class="dropdown-item"
                        :class="{ 'active': settings.tmdbProxy === proxy.id }"
                        @click="selectSettingsProxy(proxy.id)"
                      >
                        {{ proxy.name }}
                        <svg v-if="settings.tmdbProxy === proxy.id" class="dropdown-check" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </button>
                    </div>
                  </Transition>
                </div>
                <button 
                  class="btn btn-check" 
                  :class="{ 'checking': checkingProxy }"
                  @click="checkProxy"
                  :disabled="checkingProxy"
                >
                  <span v-if="checkingProxy" class="spinner-small"></span>
                  <span v-else>Проверить</span>
                </button>
              </div>
              <div class="settings-status" v-if="proxyStatus !== null">
                <span :class="proxyStatus ? 'status-online' : 'status-offline'">
                  {{ proxyStatus ? '✓ Прокси доступен' : '✗ Прокси недоступен' }}
                </span>
              </div>
              <span class="settings-hint">Прокси для получения информации о фильмах</span>
            </div>
          </div>
          
          <!-- About Section -->
          <div class="settings-section about-section">
            <h3 class="settings-section-title">О приложении</h3>
            <div class="about-content">
              <div class="about-logo">
                <img src="./assets/Knit.png" alt="Knit" class="about-icon" />
              </div>
              <div class="about-info">
                <h4 class="about-name">Knit</h4>
                <p class="about-version">v1.0.0</p>
                <p class="about-description">Элегантное приложение для поиска торрентов фильмов</p>
                <div class="about-links">
                  <a href="https://github.com/Cheviiot" target="_blank" class="about-link">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    Cheviiot
                  </a>
                  <a href="https://github.com/Cheviiot/Knit" target="_blank" class="about-link">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z"/>
                      <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
                      <line x1="12" y1="22.08" x2="12" y2="12"/>
                    </svg>
                    GitHub
                  </a>
                </div>
              </div>
            </div>
          </div>
          
          <div class="settings-footer">
            <button class="btn btn-secondary" @click="showSettings = false">Отмена</button>
            <button class="btn btn-primary" @click="saveSettings">Сохранить</button>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- Toast Notification -->
    <Transition name="toast">
      <div v-if="toast" class="toast" :class="'toast-' + toast.type">
        {{ toast.message }}
      </div>
    </Transition>
  </div>
</template>

<script>
import { GetPublicServers, SearchTorrents, SearchTMDB, SearchTMDBTV, GetPopularMovies, GetPopularTVShows, SearchWithMovie, SearchWithTVShow, GetSettings, SaveSettings, GetTMDBProxies, CheckServerByID, CheckTMDBProxyByID, GetCurrentImageBase, AddMagnet, AddTorrentURL, GetDownloads, RemoveDownload, OpenDownloadFolder, OpenFDM, CopyMagnetToClipboard } from '../bindings/knit/app.js'
import { Application, Window, Clipboard } from '@wailsio/runtime'

// Window helpers
const WindowMinimise = () => Window.Minimise()
const WindowToggleMaximise = () => Window.ToggleMaximise()
const Quit = () => Application.Quit()

export default {
  name: 'App',
  data() {
    return {
      searchQuery: '',
      contentType: 'movies', // 'movies' or 'tv'
      movies: [],
      tvShows: [],
      servers: [],
      tmdbProxies: [],
      selectedServer: 'jacred_xyz',
      selectedMovie: null,
      selectedTVShow: null,
      torrents: [],
      downloads: [],
      showDownloads: false,
      downloadsPollInterval: null,
      torrentsLoading: false,
      torrentsError: null,
      loading: false,
      loadingServers: false,
      loadingMore: false,
      isSearchResult: false,
      showSettings: false,
      currentPage: 1,
      hasMore: true,
      settings: {
        selectedServer: 'jacred_xyz',
        tmdbProxy: 'apn_render',
        theme: 'dark'
      },
      toast: null,
      checkingServer: false,
      checkingProxy: false,
      serverStatus: null,
      proxyStatus: null,
      imageBase: 'https://image.tmdb.org',
      openDropdown: null,
      scrollThrottleTimer: null,
      resizeDebounceTimer: null,
      imageRetries: {} // Track retry counts for images
    }
  },
  computed: {
    currentServerName() {
      const server = this.servers.find(s => s.id === this.selectedServer)
      return server ? server.name : 'Сервер'
    },
    settingsServerName() {
      const server = this.servers.find(s => s.id === this.settings.selectedServer)
      return server ? server.name : 'Выбрать сервер'
    },
    settingsProxyName() {
      const proxy = this.tmdbProxies.find(p => p.id === this.settings.tmdbProxy)
      return proxy ? proxy.name : 'Выбрать прокси'
    },
    isModalOpen() {
      return this.selectedMovie || this.selectedTVShow || this.showSettings || this.showDownloads
    }
  },
  watch: {
    isModalOpen(open) {
      document.body.style.overflow = open ? 'hidden' : ''
    }
  },
  async mounted() {
    // Disable browser-like behavior
    document.addEventListener('contextmenu', e => e.preventDefault())
    document.addEventListener('selectstart', e => e.preventDefault())
    document.addEventListener('dragstart', e => e.preventDefault())
    
    await this.loadSettings()
    await this.refreshServers()
    await this.loadTMDBProxies()
    await this.updateImageBase()
    await this.loadPopularMovies()
    await this.loadDownloads()
    
    // Poll downloads every 2 seconds
    this.downloadsPollInterval = setInterval(() => {
      if (this.downloads.length > 0) {
        this.loadDownloads()
      }
    }, 2000)
    
    // Add scroll listener for infinite scroll (with throttle)
    window.addEventListener('scroll', this.throttledHandleScroll, { passive: true })
    window.addEventListener('resize', this.debouncedCheckNeedMoreContent, { passive: true })
    
    // Close dropdown on click outside
    document.addEventListener('click', this.closeDropdownOnClickOutside)
  },
  beforeUnmount() {
    if (this.downloadsPollInterval) {
      clearInterval(this.downloadsPollInterval)
    }
    window.removeEventListener('scroll', this.throttledHandleScroll)
    window.removeEventListener('resize', this.debouncedCheckNeedMoreContent)
    if (this.scrollThrottleTimer) clearTimeout(this.scrollThrottleTimer)
    if (this.resizeDebounceTimer) clearTimeout(this.resizeDebounceTimer)
    document.removeEventListener('click', this.closeDropdownOnClickOutside)
  },
  methods: {
    toggleDropdown(name) {
      this.openDropdown = this.openDropdown === name ? null : name
    },
    
    closeDropdownOnClickOutside(e) {
      if (!e.target.closest('.dropdown')) {
        this.openDropdown = null
      }
    },
    
    selectServer(id) {
      this.selectedServer = id
      this.openDropdown = null
    },
    
    selectSettingsServer(id) {
      this.settings.selectedServer = id
      this.serverStatus = null
      this.openDropdown = null
    },
    
    selectSettingsProxy(id) {
      this.settings.tmdbProxy = id
      this.proxyStatus = null
      this.openDropdown = null
    },
    
    async loadDownloads() {
      try {
        this.downloads = await GetDownloads() || []
      } catch (e) {
        // Ignore errors during polling
      }
    },
    
    async downloadTorrent(torrent) {
      const name = torrent.title || 'Торрент'
      
      try {
        this.showToast('Отправляю в FDM...', 'info')
        
        let result
        if (torrent.magnetUri) {
          result = await AddMagnet(torrent.magnetUri, name)
        } else if (torrent.link) {
          result = await AddTorrentURL(torrent.link, name)
        } else {
          this.showToast('Ссылка недоступна', 'error')
          return
        }
        
        await this.loadDownloads()
        this.showDownloads = true
        if (result && result.client) {
          this.showToast(`Отправлено в ${result.client}!`, 'success')
        } else {
          this.showToast('Отправлено в FDM!', 'success')
        }
      } catch (e) {
        this.showToast(`Ошибка: ${e.message || e}`, 'error')
      }
    },
    
    async removeDownload(id) {
      try {
        await RemoveDownload(id)
        await this.loadDownloads()
        this.showToast('Загрузка удалена', 'success')
      } catch {
        this.showToast('Ошибка удаления', 'error')
      }
    },
    
    async openDownloadsFolder() {
      try {
        await OpenDownloadFolder()
      } catch {
        this.showToast('Не удалось открыть папку', 'error')
      }
    },
    
    async openFDM() {
      try {
        await OpenFDM()
      } catch {
        this.showToast('Не удалось открыть FDM', 'error')
      }
    },
    
    getStatusText(status) {
      const texts = {
        'downloading': 'Загрузка',
        'seeding': 'Раздача',
        'paused': 'Пауза',
        'completed': 'Завершено',
        'error': 'Ошибка',
        'sending': 'Отправка...',
        'sent_to_fdm': '✓ В FDM',
        'sent_to_qbittorrent': '✓ В qBittorrent',
        'sent_to_transmission': '✓ В Transmission',
        'sent_to_deluge': '✓ В Deluge',
        'sent_to_ktorrent': '✓ В KTorrent',
        'added': 'Добавлен',
        'opened': '✓ Открыт',
        'downloading_torrent': 'Скачивание .torrent...'
      }
      return texts[status] || status
    },
    
    async updateImageBase() {
      // No longer needed - using imagetmdb.com proxy
    },
    
    getImageUrl(path, size = 'w500') {
      if (!path) return ''
      // Use local proxy to avoid WebKit CSP blocking external images
      return `/tmdb-image/${size}${path}`
    },
    
    // Handle image load error with retry
    handleImageError(event, path, size = 'w500') {
      const key = `${path}_${size}`
      const maxRetries = 3
      const retryDelay = 1000 // 1 second
      
      if (!this.imageRetries[key]) {
        this.imageRetries[key] = 0
      }
      
      this.imageRetries[key]++
      
      if (this.imageRetries[key] <= maxRetries) {
        // Retry loading via backend
        setTimeout(() => {
          const newUrl = this.getImageUrl(path, size) + `?retry=${this.imageRetries[key]}`
          event.target.src = newUrl
        }, retryDelay * this.imageRetries[key])
      } else {
        // Max retries reached, hide the image or show placeholder
        event.target.style.display = 'none'
        // Show the no-poster placeholder by triggering parent's fallback
        const posterDiv = event.target.closest('.movie-poster')
        if (posterDiv) {
          const placeholder = posterDiv.querySelector('.no-poster')
          if (placeholder) {
            placeholder.style.display = 'flex'
          }
        }
      }
    },
    
    // Reset image retries when content changes
    resetImageRetries() {
      this.imageRetries = {}
    },
    
    async loadSettings() {
      try {
        const settings = await GetSettings()
        if (settings) {
          this.settings = { ...this.settings, ...settings }
          this.selectedServer = this.settings.selectedServer || 'jacred_xyz'
        }
      } catch {
        // Use default settings
      }
    },
    
    async saveSettings() {
      try {
        const oldTMDBProxy = this.settings.tmdbProxy
        await SaveSettings(this.settings)
        this.selectedServer = this.settings.selectedServer
        this.showSettings = false
        this.serverStatus = null
        this.proxyStatus = null
        this.showToast('Настройки сохранены', 'success')
        
        // Reload movies if TMDB proxy changed
        if (oldTMDBProxy !== this.settings.tmdbProxy) {
          await this.updateImageBase()
          await this.loadPopularMovies()
        }
      } catch (e) {
        this.showToast('Ошибка сохранения', 'error')
      }
    },
    
    async checkServer() {
      this.checkingServer = true
      this.serverStatus = null
      try {
        this.serverStatus = await CheckServerByID(this.settings.selectedServer)
      } catch (e) {
        this.serverStatus = false
      }
      this.checkingServer = false
    },
    
    async checkProxy() {
      this.checkingProxy = true
      this.proxyStatus = null
      try {
        this.proxyStatus = await CheckTMDBProxyByID(this.settings.tmdbProxy)
      } catch (e) {
        this.proxyStatus = false
      }
      this.checkingProxy = false
    },
    
    async refreshServers() {
      this.loadingServers = true
      try {
        this.servers = await GetPublicServers()
      } catch {
        // Keep existing servers on error
      }
      this.loadingServers = false
    },
    
    async loadTMDBProxies() {
      try {
        this.tmdbProxies = await GetTMDBProxies()
      } catch {
        // Keep existing proxies on error
      }
    },
    
    async loadPopularMovies(reset = true) {
      if (reset) {
        this.loading = true
        this.currentPage = 1
        this.hasMore = true
        this.movies = []
      } else {
        this.loadingMore = true
      }
      
      this.isSearchResult = false
      
      try {
        const response = await GetPopularMovies(this.currentPage)
        if (response && response.results) {
          if (reset) {
            this.movies = response.results
          } else {
            this.movies = [...this.movies, ...response.results]
          }
          
          // Check if there are more pages
          this.hasMore = this.currentPage < response.total_pages && response.results.length > 0
          this.currentPage++
        }
      } catch {
        this.hasMore = false
      }
      
      this.loading = false
      this.loadingMore = false
      
      // Check if we need to load more content (e.g., on large screens)
      this.$nextTick(() => {
        this.checkNeedMoreContent()
      })
    },
    
    checkNeedMoreContent() {
      // If content doesn't fill the screen, load more
      if (this.loadingMore || !this.hasMore || this.isSearchResult || this.loading) {
        return
      }
      
      const scrollHeight = document.documentElement.scrollHeight
      const clientHeight = window.innerHeight
      
      // If there's no scrollbar (content fits on screen), load more
      if (scrollHeight <= clientHeight + 100) {
        if (this.contentType === 'movies') {
          this.loadPopularMovies(false)
        } else {
          this.loadPopularTVShows(false)
        }
      }
    },
    
    // Debounced version for resize events
    debouncedCheckNeedMoreContent() {
      if (this.resizeDebounceTimer) clearTimeout(this.resizeDebounceTimer)
      this.resizeDebounceTimer = setTimeout(() => {
        this.checkNeedMoreContent()
      }, 150)
    },
    
    // Throttled version for scroll events
    throttledHandleScroll() {
      if (this.scrollThrottleTimer) return
      this.scrollThrottleTimer = setTimeout(() => {
        this.handleScroll()
        this.scrollThrottleTimer = null
      }, 100)
    },
    
    handleScroll() {
      // Don't load more if already loading, no more data, or in search mode
      if (this.loadingMore || !this.hasMore || this.isSearchResult || this.loading) {
        return
      }
      
      const scrollTop = window.pageYOffset || document.documentElement.scrollTop
      const scrollHeight = document.documentElement.scrollHeight
      const clientHeight = window.innerHeight
      
      // Load more when user is 300px from bottom
      if (scrollTop + clientHeight >= scrollHeight - 300) {
        if (this.contentType === 'movies') {
          this.loadPopularMovies(false)
        } else {
          this.loadPopularTVShows(false)
        }
      }
    },
    
    async searchMovies() {
      if (!this.searchQuery.trim()) {
        await this.loadPopularMovies()
        return
      }
      
      this.loading = true
      this.isSearchResult = true
      try {
        const response = await SearchTMDB(this.searchQuery)
        if (response && response.results) {
          this.movies = response.results
        }
      } catch {
        this.showToast('Ошибка поиска', 'error')
      }
      this.loading = false
    },
    
    clearSearch() {
      this.searchQuery = ''
      if (this.contentType === 'movies') {
        this.loadPopularMovies()
      } else {
        this.loadPopularTVShows()
      }
    },
    
    selectMovie(movie) {
      this.selectedMovie = movie
      this.torrents = []
      this.torrentsError = null
      this.searchTorrents()
    },
    
    closeModal() {
      this.selectedMovie = null
      this.torrents = []
      this.torrentsError = null
    },
    
    async searchTorrents() {
      if (!this.selectedMovie) return
      
      this.torrentsLoading = true
      this.torrentsError = null
      this.torrents = []
      
      try {
        const results = await SearchWithMovie(this.selectedMovie, this.selectedServer)
        this.torrents = results || []
        
        if (this.torrents.length === 0) {
          this.torrentsError = 'Торренты не найдены. Попробуйте другой сервер.'
        }
      } catch (e) {
        this.torrentsError = `Ошибка: ${e.message || e}`
      }
      
      this.torrentsLoading = false
    },
    
    // TV Shows methods
    selectTVShow(show) {
      this.selectedTVShow = show
      this.torrents = []
      this.torrentsError = null
      this.searchTVTorrents()
    },
    
    closeTVModal() {
      this.selectedTVShow = null
      this.torrents = []
      this.torrentsError = null
    },
    
    async searchTVTorrents() {
      if (!this.selectedTVShow) return
      
      this.torrentsLoading = true
      this.torrentsError = null
      this.torrents = []
      
      try {
        const results = await SearchWithTVShow(this.selectedTVShow, this.selectedServer)
        this.torrents = results || []
        
        if (this.torrents.length === 0) {
          this.torrentsError = 'Торренты не найдены. Попробуйте другой сервер.'
        }
      } catch (e) {
        this.torrentsError = `Ошибка: ${e.message || e}`
      }
      
      this.torrentsLoading = false
    },
    
    async loadPopularTVShows(reset = true) {
      if (reset) {
        this.loading = true
        this.currentPage = 1
        this.hasMore = true
        this.tvShows = []
      } else {
        this.loadingMore = true
      }
      
      this.isSearchResult = false
      
      try {
        const response = await GetPopularTVShows(this.currentPage)
        if (response && response.results) {
          if (reset) {
            this.tvShows = response.results
          } else {
            this.tvShows = [...this.tvShows, ...response.results]
          }
          
          // Check if there are more pages
          this.hasMore = this.currentPage < response.total_pages && response.results.length > 0
          this.currentPage++
        }
      } catch {
        this.hasMore = false
      }
      
      this.loading = false
      this.loadingMore = false
      
      // Check if we need to load more content (e.g., on large screens)
      this.$nextTick(() => {
        this.checkNeedMoreContent()
      })
    },
    
    async searchTVShows() {
      if (!this.searchQuery.trim()) {
        await this.loadPopularTVShows()
        return
      }
      
      this.loading = true
      this.isSearchResult = true
      try {
        const response = await SearchTMDBTV(this.searchQuery)
        if (response && response.results) {
          this.tvShows = response.results
        }
      } catch {
        this.showToast('Ошибка поиска', 'error')
      }
      this.loading = false
    },
    
    switchContentType(type) {
      if (this.contentType === type) return
      this.contentType = type
      this.searchQuery = ''
      this.isSearchResult = false
      this.currentPage = 1
      this.hasMore = true
      
      if (type === 'movies') {
        this.tvShows = []
        this.loadPopularMovies()
      } else {
        this.movies = []
        this.loadPopularTVShows()
      }
    },
    
    async searchContent() {
      if (this.contentType === 'movies') {
        await this.searchMovies()
      } else {
        await this.searchTVShows()
      }
    },
    
    async copyMagnet(magnetUri) {
      try {
        await CopyMagnetToClipboard(magnetUri)
        this.showToast('Magnet-ссылка скопирована', 'success')
      } catch {
        this.showToast('Не удалось скопировать', 'error')
      }
    },
    
    getYear(date) {
      return date ? date.split('-')[0] : ''
    },
    
    formatDate(date) {
      if (!date) return ''
      const d = new Date(date)
      return d.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' })
    },
    
    showToast(message, type = 'info') {
      this.toast = { message, type }
      setTimeout(() => {
        this.toast = null
      }, 3000)
    },
    
    minimizeWindow() {
      WindowMinimise()
    },
    
    toggleMaximize() {
      WindowToggleMaximise()
    },
    
    closeWindow() {
      Quit()
    }
  }
}
</script>

<style>
/* Disable all browser-like behavior */
*, *::before, *::after {
  scrollbar-width: none;
  -ms-overflow-style: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  -webkit-touch-callout: none;
  -webkit-tap-highlight-color: transparent;
}

/* Scroll optimization */
html {
  scroll-behavior: auto;
}

body {
  overscroll-behavior: none;
  -webkit-touch-callout: none;
}

*::-webkit-scrollbar {
  display: none;
}

/* Disable image/link dragging */
img, a {
  -webkit-user-drag: none;
  pointer-events: none;
}

img {
  pointer-events: none;
}

/* Disable context menu */
body {
  -webkit-touch-callout: none;
}
</style>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Titlebar */
.titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 0 0 14px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.titlebar-center {
  flex: 1;
  max-width: 420px;
  padding: 0 20px;
}

.titlebar-right {
  display: flex;
  align-items: center;
}

/* Brand */
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  width: 22px;
  height: 22px;
  border-radius: var(--radius-xs);
  object-fit: contain;
}

.brand-name {
  font-size: 14px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-hover) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Custom Dropdown */
.dropdown {
  position: relative;
}

.dropdown-full {
  flex: 1;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
  width: 100%;
}

.dropdown-trigger:hover {
  border-color: var(--border-light);
}

.dropdown-trigger:focus {
  outline: none;
  border-color: var(--accent);
}

.server-trigger {
  padding: 4px 8px;
  border-radius: 14px;
  width: auto;
}

.dropdown-input {
  justify-content: space-between;
  padding: 8px 12px;
}

.dropdown-value {
  flex: 1;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dropdown-arrow {
  flex-shrink: 0;
  color: var(--text-muted);
  transition: transform 0.2s;
}

.dropdown-trigger:hover .dropdown-arrow {
  color: var(--text-secondary);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  min-width: 180px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  padding: 6px;
  z-index: 1000;
  max-height: 240px;
  overflow-y: auto;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.1s;
  text-align: left;
}

.dropdown-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.dropdown-item.active {
  color: var(--accent);
}

.dropdown-item-dot {
  width: 6px;
  height: 6px;
  background: var(--text-muted);
  border-radius: 50%;
  flex-shrink: 0;
}

.dropdown-item-dot.online {
  background: #22c55e;
}

.dropdown-check {
  margin-left: auto;
  color: var(--accent);
  flex-shrink: 0;
}

.server-dot {
  width: 6px;
  height: 6px;
  background: #22c55e;
  border-radius: 50%;
  flex-shrink: 0;
}

/* Dropdown Transition */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* Content Type Toggle */
.content-type-toggle {
  display: flex;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 2px;
  gap: 2px;
  margin-right: 8px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 28px;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: var(--transition);
}

.toggle-btn:hover {
  color: var(--text-primary);
  background: var(--bg-tertiary);
}

.toggle-btn.active {
  background: var(--accent);
  color: white;
}

.toggle-btn svg {
  flex-shrink: 0;
}

/* Search Box */
.search-box {
  display: flex;
  align-items: center;
  height: 36px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-full);
  transition: var(--transition);
}

.search-box:focus-within {
  border-color: var(--accent);
  background: var(--bg-tertiary);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.search-icon {
  margin-left: 14px;
  color: var(--text-muted);
  flex-shrink: 0;
  transition: var(--transition-fast);
}

.search-box:focus-within .search-icon {
  color: var(--accent);
}

.search-input {
  flex: 1;
  height: 100%;
  padding: 0 12px 0 10px;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
}

.search-input::placeholder {
  color: var(--text-muted);
  transition: var(--transition-fast);
}

.search-box:focus-within .search-input::placeholder {
  opacity: 0.5;
}

.search-clear {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  margin-right: 4px;
  background: transparent;
  border: none;
  border-radius: var(--radius-full);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.search-clear:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  transform: rotate(90deg);
}

/* Titlebar Button */
.titlebar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  margin-right: 4px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.titlebar-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

/* Titlebar Divider */
.titlebar-divider {
  width: 1px;
  height: 18px;
  background: var(--border);
  margin: 0 6px;
}

/* Window Controls */
.window-controls {
  display: flex;
  height: 100%;
  gap: 4px;
  padding: 0 8px;
  align-items: center;
}

.window-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: var(--transition-fast);
}

.window-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.window-btn:active {
  transform: scale(0.92);
}

.window-btn.window-close:hover {
  background: rgba(239, 68, 68, 0.15);
  color: var(--error);
}

/* Main */
.main {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.spin {
  animation: spin 1s linear infinite;
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-secondary);
  gap: 16px;
}

/* Loading More */
.loading-more {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--text-secondary);
  gap: 12px;
}

.loading-more .spinner {
  width: 32px;
  height: 32px;
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-secondary);
  text-align: center;
}

.empty-state svg {
  color: var(--text-muted);
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

/* Movies Section */
.movies-section {
  animation: fadeIn 0.3s ease;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  margin-bottom: 20px;
}

.section-title .count {
  font-weight: 400;
  color: var(--text-muted);
}

/* Main Content */
.main {
  margin-top: 0;
}

/* Movies Grid */
.movies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 24px;
  contain: layout style;
}

.movie-card {
  cursor: pointer;
  transition: var(--transition);
  animation: fadeIn 0.4s ease;
  contain: layout style paint;
  will-change: transform;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-card);
  border: 1px solid transparent;
}

.movie-card:hover {
  transform: translateY(-6px);
  border-color: var(--border-light);
  box-shadow: var(--shadow-md);
}

.movie-card:hover .movie-poster {
  box-shadow: none;
}

.movie-poster {
  position: relative;
  aspect-ratio: 2/3;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  overflow: hidden;
  background: linear-gradient(135deg, var(--bg-tertiary) 0%, var(--bg-secondary) 100%);
  transition: var(--transition);
}

.movie-poster::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, var(--bg-card) 0%, transparent 40%);
  opacity: 0;
  transition: var(--transition);
}

.movie-card:hover .movie-poster::after {
  opacity: 1;
}

.movie-poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  content-visibility: auto;
}

.poster-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--bg-tertiary) 0%, var(--bg-secondary) 100%);
}

.spinner-small {
  width: 24px;
  height: 24px;
  border: 2px solid var(--bg-tertiary);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.no-poster {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: var(--text-muted);
  background: linear-gradient(135deg, var(--bg-tertiary) 0%, var(--bg-secondary) 100%);
}

.movie-rating {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 8px;
  background: var(--glass-bg);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-radius: var(--radius-full);
  font-size: 11px;
  font-weight: 700;
  color: var(--accent);
  border: 1px solid var(--border-accent);
}

.movie-info {
  padding: 12px;
}

.movie-title {
  font-size: 13px;
  font-weight: 600;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: var(--transition-fast);
  color: var(--text-primary);
}

.movie-card:hover .movie-title {
  color: var(--accent);
}

.movie-year {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .movie-modal,
.modal-leave-to .movie-modal,
.modal-enter-from .settings-modal,
.modal-leave-to .settings-modal {
  transform: scale(0.9) translateY(20px);
}

.movie-modal {
  position: relative;
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  background: var(--glass-bg);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-radius: var(--radius-xl);
  overflow: hidden;
  overflow-y: auto;
  box-shadow: var(--shadow-lg), 0 0 80px rgba(0, 0, 0, 0.4);
  border: 1px solid var(--glass-border);
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 10;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-full);
  padding: 0;
  color: var(--text-secondary);
  cursor: pointer;
  transition: var(--transition);
}

.modal-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
  transform: rotate(90deg);
  border-color: var(--border-light);
}

.modal-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 400px;
  overflow: hidden;
  z-index: 0;
}

.modal-backdrop img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.backdrop-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom, 
    rgba(30, 30, 30, 0.3) 0%, 
    rgba(30, 30, 30, 0.7) 50%,
    var(--bg-secondary) 100%
  );
}

.modal-content {
  display: flex;
  gap: 24px;
  padding: 24px;
  padding-top: 40px;
  position: relative;
  z-index: 1;
  align-items: flex-start;
  min-height: 320px;
}

.modal-poster {
  flex-shrink: 0;
  width: 240px;
  height: fit-content;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.modal-poster img {
  width: 100%;
  display: block;
}

.modal-info {
  flex: 1;
  min-width: 0;
}

.modal-title {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 4px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}

.modal-original {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 16px;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.5);
}

.modal-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--text-secondary);
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.5);
}

.meta-item.rating {
  color: #fbbf24;
}

.modal-overview {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-bottom: 16px;
}

/* Torrents Section */
.torrents-section {
  padding: 0 20px 20px;
}

.torrents-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.torrents-title .count {
  font-weight: 400;
  color: var(--text-muted);
}

.torrents-loading {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-secondary);
  padding: 20px 0;
}

.torrents-error {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--error);
  padding: 16px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: var(--radius-md);
}

.torrents-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.torrent-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  transition: var(--transition);
  border: 1px solid transparent;
}

.torrent-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-light);
  transform: translateX(4px);
}

.torrent-info {
  flex: 1;
  min-width: 0;
}

.torrent-title {
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  margin-bottom: 6px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.torrent-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-muted);
}

.torrent-stats {
  display: flex;
  gap: 12px;
}

.seeders,
.leechers {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 6px;
}

.seeders {
  color: var(--success);
  background: rgba(16, 185, 129, 0.1);
}

.leechers {
  color: var(--error);
  background: rgba(239, 68, 68, 0.1);
}

.torrent-actions {
  display: flex;
  gap: 8px;
}

/* Settings Modal */
.settings-modal {
  width: 100%;
  max-width: 520px;
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--glass-border);
  overflow: hidden;
  box-shadow: var(--shadow-lg), 0 0 0 1px rgba(255, 255, 255, 0.03);
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid var(--border);
}

.settings-header h2 {
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--text-primary) 0%, var(--text-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.settings-content {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.settings-section {
  margin-bottom: 20px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-muted);
  margin-bottom: 12px;
}

.settings-group {
  margin-bottom: 20px;
}

.settings-group:last-child {
  margin-bottom: 0;
}

.settings-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-secondary);
}

.settings-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.settings-row .input {
  flex: 1;
}

.btn-check {
  padding: 10px 16px;
  min-width: 100px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  color: var(--text-primary);
  border-radius: var(--radius-sm);
  font-size: 13px;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-check:hover:not(:disabled) {
  background: var(--bg-secondary);
  border-color: var(--accent);
}

.btn-check:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-check.checking {
  background: var(--accent);
  color: #000;
}

.spinner-small {
  width: 14px;
  height: 14px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.settings-status {
  margin-top: 8px;
  font-size: 13px;
}

.status-online {
  color: #4ade80;
}

.status-offline {
  color: #f87171;
}

.settings-hint {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 6px;
}

.settings-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 8px;
}

.settings-desc a {
  color: var(--accent);
  text-decoration: none;
}

.settings-desc a:hover {
  text-decoration: underline;
}

.settings-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-tertiary);
  border-top: 1px solid var(--border);
}

/* About Section */
.about-section {
  border-top: 1px solid var(--border);
  padding-top: 16px;
  margin-top: 8px;
}

.about-content {
  display: flex;
  gap: 16px;
  align-items: center;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.about-logo {
  flex-shrink: 0;
}

.about-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
}

.about-info {
  flex: 1;
}

.about-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.about-version {
  font-size: 12px;
  color: var(--text-muted);
  margin: 2px 0 8px 0;
}

.about-description {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0 0 12px 0;
}

.about-links {
  display: flex;
  gap: 12px;
}

.about-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 13px;
  transition: all 0.2s;
}

.about-link:hover {
  background: var(--accent);
  color: #000;
  border-color: var(--accent);
}

.about-link svg {
  flex-shrink: 0;
}

/* Toast */
.toast {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  z-index: 2000;
  box-shadow: var(--shadow-lg);
}

.toast-success {
  border-color: var(--success);
  color: var(--success);
}

.toast-error {
  border-color: var(--error);
  color: var(--error);
}

.toast-warning {
  border-color: var(--warning);
  color: var(--warning);
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}

/* Downloads Drawer */
.downloads-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 1000;
}

.downloads-drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 380px;
  max-width: 100%;
  background: var(--bg-primary);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  box-shadow: -8px 0 32px rgba(0, 0, 0, 0.3);
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  border-bottom: 1px solid var(--border);
  background: var(--bg-secondary);
}

.drawer-title {
  display: flex;
  align-items: center;
  gap: 14px;
}

.drawer-icon {
  width: 42px;
  height: 42px;
  background: linear-gradient(135deg, var(--accent), #e69500);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
}

.drawer-title h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.drawer-subtitle {
  font-size: 12px;
  color: var(--text-muted);
}

.drawer-actions {
  display: flex;
  gap: 8px;
}

.drawer-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: var(--bg-tertiary);
  border-radius: 10px;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.drawer-btn:hover {
  background: var(--border);
  color: var(--text-primary);
}

.drawer-close:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.drawer-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.drawer-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: var(--text-muted);
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: var(--bg-secondary);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  opacity: 0.6;
}

.drawer-empty p {
  font-size: 15px;
  font-weight: 500;
  margin: 0 0 4px;
  color: var(--text-secondary);
}

.drawer-empty span {
  font-size: 13px;
}

.drawer-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.dl-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border);
  transition: all 0.2s;
}

.dl-card:hover {
  border-color: var(--accent);
  background: var(--bg-tertiary);
}

.dl-status-indicator {
  width: 4px;
  height: 36px;
  border-radius: 2px;
  background: var(--text-muted);
  flex-shrink: 0;
}

.dl-status-indicator.indicator-fdm,
.dl-status-indicator.indicator-qbittorrent,
.dl-status-indicator.indicator-transmission,
.dl-status-indicator.indicator-deluge,
.dl-status-indicator.indicator-ktorrent,
.dl-status-indicator.indicator-completed,
.dl-status-indicator.indicator-seeding {
  background: #22c55e;
}

.dl-status-indicator.indicator-downloading,
.dl-status-indicator.indicator-sending {
  background: var(--accent);
}

.dl-status-indicator.indicator-error {
  background: #ef4444;
}

.dl-content {
  flex: 1;
  min-width: 0;
}

.dl-name {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.dl-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-muted);
}

.dl-status {
  font-weight: 500;
}

.dl-status.status-fdm,
.dl-status.status-qbittorrent,
.dl-status.status-transmission,
.dl-status.status-deluge,
.dl-status.status-ktorrent,
.dl-status.status-completed,
.dl-status.status-seeding {
  color: #22c55e;
}

.dl-status.status-downloading,
.dl-status.status-sending {
  color: var(--accent);
}

.dl-status.status-error {
  color: #ef4444;
}

.dl-actions {
  display: flex;
  gap: 6px;
}

.dl-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: var(--bg-tertiary);
  border-radius: 8px;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.dl-btn:hover {
  background: var(--accent);
  color: #000;
}

.dl-btn-remove:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

/* Downloads FAB */
.downloads-fab {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 44px;
  height: 44px;
  border: none;
  background: linear-gradient(135deg, var(--accent) 0%, #e6a800 100%);
  border-radius: 12px;
  color: #000;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(245, 197, 24, 0.35);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 100;
}

.downloads-fab:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 24px rgba(245, 197, 24, 0.45);
}

.downloads-fab.fab-active {
  transform: scale(0.95);
  animation: none;
}

.fab-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  border-radius: 9px;
  font-size: 10px;
  font-weight: 700;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(239, 68, 68, 0.4);
}

/* Drawer Transition */
.drawer-enter-active,
.drawer-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.drawer-enter-active .downloads-drawer,
.drawer-leave-active .downloads-drawer {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.drawer-enter-from,
.drawer-leave-to {
  background: rgba(0, 0, 0, 0);
}

.drawer-enter-from .downloads-drawer,
.drawer-leave-to .downloads-drawer {
  transform: translateX(100%);
}

.btn-accent {
  background: var(--accent);
  color: #000;
}

.btn-accent:hover {
  background: #f0b429;
}

.btn-danger {
  color: #ef4444;
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 0.1);
}

/* Slide up animation */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}

/* Responsive */
@media (max-width: 768px) {
  .header {
    flex-wrap: wrap;
    gap: 12px;
  }
  
  .header-center {
    order: 3;
    flex: 100%;
    max-width: none;
  }
  
  .modal-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
    margin-top: -60px;
  }
  
  .modal-poster {
    width: 140px;
  }
  
  .modal-meta {
    justify-content: center;
  }
  
  .movies-grid {
    grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
    gap: 16px;
  }
  
  .torrent-item {
    flex-wrap: wrap;
  }
  
  .torrent-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

/* Масштабирование модальных окон */
@media (max-height: 700px) {
  .movie-modal {
    max-height: 95vh;
  }
  
  .modal-backdrop {
    height: 300px;
  }
  
  .modal-content {
    padding: 16px;
    padding-top: 24px;
    min-height: 260px;
  }
  
  .modal-poster {
    width: 180px;
  }
  
  .modal-title {
    font-size: 20px;
  }
  
  .modal-overview {
    font-size: 13px;
    -webkit-line-clamp: 3;
    line-clamp: 3;
  }
}

@media (max-height: 550px) {
  .modal-backdrop {
    height: 220px;
  }
  
  .modal-content {
    padding: 12px;
    padding-top: 16px;
    gap: 16px;
    min-height: 200px;
  }
  
  .modal-poster {
    width: 140px;
  }
  
  .modal-title {
    font-size: 18px;
  }
  
  .modal-overview {
    font-size: 12px;
    -webkit-line-clamp: 2;
    line-clamp: 2;
  }
  
  .torrents-section {
    padding: 12px;
  }
  
  .torrent-item {
    padding: 8px 12px;
  }
}

@media (max-width: 600px) {
  .movie-modal {
    max-width: 95vw;
    border-radius: var(--radius-lg);
  }
  
  .modal-backdrop {
    height: 250px;
  }
  
  .modal-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 16px;
    padding-top: 20px;
    gap: 16px;
  }
  
  .modal-poster {
    width: 160px;
  }
  
  .modal-meta {
    justify-content: center;
  }
  
  .settings-modal {
    max-width: 95vw;
    margin: 16px;
  }
}
</style>
