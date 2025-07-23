// ===== BACKEND CONFIGURATION =====
// Change this URL to point to your backend server
const BACKEND_CONFIG = {
  // Empty string = same server (frontend and backend on same domain)
  // Example: 'http://localhost:8001' for different server
  // Example: 'https://api.yourbackend.com' for production
  BASE_URL: 'http://localhost:8082',
  
  // API endpoints (don't change these unless your backend uses different paths)
  ENDPOINTS: {
    LOGIN: '/login',
    SIGNUP: '/signup',
    GET_SYLLABUS: '/get_syllabus',
    SAVE_SYLLABUS: '/save_syllabus',
    DELETE_SYLLABUS: '/delete_syllabus',
    GENERATE_TASK: '/generate_task',
    SUBMIT_CODE: '/submit_code',
    NOTIFICATION_SETTINGS: '/notification-settings',
    GET_STATS: '/get_stats'
  }
};

// Helper function to build API URLs
const apiUrl = (endpoint) => {
  return BACKEND_CONFIG.BASE_URL + endpoint;
}; 