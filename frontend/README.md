# Python Practice - Frontend

A modern, interactive coding platform for Python practice with real-time feedback.

## 🚀 Quick Start

1. **Start the frontend server:**
   ```bash
   python3 -m http.server 8000
   ```

2. **Open in browser:**
   ```
   http://localhost:8000
   ```

## ⚙️ Backend Configuration

### **Where to Configure Backend URL:**

Edit `config.js` file and change the `BASE_URL`:

```javascript
const BACKEND_CONFIG = {
  // Change this to your backend server URL
  BASE_URL: '', // Empty = same server, or 'http://localhost:8001' for different server
  // ...
};
```

### **Configuration Examples:**

| Backend Setup | BASE_URL | Description |
|---------------|----------|-------------|
| **Same Server** | `''` | Frontend and backend on same domain |
| **Local Different Port** | `'http://localhost:8001'` | Backend on port 8001 |
| **Production API** | `'https://api.yourbackend.com'` | Remote backend server |

### **Backend API Endpoints Required:**

Your backend must implement these endpoints:

#### **Authentication:**
- `POST /login` - User login (supports both regular users and admin users)
- `POST /signup` - User registration

#### **Syllabus Management:**
- `GET /get_syllabus` - Get course topics
- `POST /save_syllabus` - Save uploaded syllabus
- `DELETE /delete_syllabus` - Clear syllabus

#### **Task Generation & Submission:**
- `POST /generate_task` - Generate coding tasks (JSON body with topic and difficulty)
- `POST /submit_code` - Submit and evaluate code

#### **User Progress:**
- `GET /get_stats` - Get user's solved task statistics

#### **Notifications:**
- `GET /notification-settings` - Get user preferences
- `POST /notification-settings` - Save preferences

## 🎯 Features

- ✅ **User Authentication** (login/signup)
- ✅ **Interactive Coding** with real-time feedback
- ✅ **Task Generation** with difficulty levels
- ✅ **Progress Tracking** and scoring (backend-based)
- ✅ **Smart Hints** system
- ✅ **Notification Settings**
- ✅ **Modern UI** with animations

## 📁 File Structure

```
frontend_study_buddy/
├── index.html          # Main HTML file
├── script.js           # Main JavaScript logic
├── style.css           # Styling
├── config.js           # Backend configuration
└── README.md           # This file
```

## 🔧 Development

- **Frontend**: Pure HTML/CSS/JavaScript
- **Backend**: Any server (Python, Node.js, etc.)
- **CORS**: Configure your backend to allow requests from frontend domain