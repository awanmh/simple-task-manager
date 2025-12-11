import { useState, useEffect, useRef } from 'react';
import api from './api';
import type { Task, LoginResponse, Subtask } from './types';
import { AxiosError } from 'axios';

// ---------------------------------------------------------
// AUTH COMPONENT
// ---------------------------------------------------------

const AuthForm = ({ onLogin }: { onLogin: (token: string) => void }) => {
  const [isRegister, setIsRegister] = useState(false);
  const [formData, setFormData] = useState({ name: '', email: '', password: '' });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (isRegister) {
        await api.post('/register', formData);
        alert("Register Berhasil! Silakan Login.");
        setIsRegister(false);
      } else {
        const res = await api.post<LoginResponse>('/login', { 
          email: formData.email, 
          password: formData.password 
        });
        onLogin(res.data.access_token);
      }
    } catch {
      alert("Gagal! Pastikan data sudah benar.");
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '50px auto', textAlign: 'center', fontFamily: 'sans-serif' }}>
      <h1>{isRegister ? 'Daftar Akun' : 'Login Task Manager'}</h1>

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        {isRegister && (
          <input
            type="text"
            placeholder="Nama Lengkap"
            required
            value={formData.name}
            onChange={e => setFormData({ ...formData, name: e.target.value })}
            style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ccc' }}
          />
        )}

        <input
          type="email"
          placeholder="Email"
          required
          value={formData.email}
          onChange={e => setFormData({ ...formData, email: e.target.value })}
          style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ccc' }}
        />

        <input
          type="password"
          placeholder="Password"
          required
          value={formData.password}
          onChange={e => setFormData({ ...formData, password: e.target.value })}
          style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ccc' }}
        />

        <button type="submit" style={{ padding: '12px', background: '#007bff', color: 'white', borderRadius: '8px', border: 'none' }}>
          {isRegister ? 'Register' : 'Login'}
        </button>
      </form>

      <p
        onClick={() => setIsRegister(!isRegister)}
        style={{ marginTop: '20px', cursor: 'pointer', color: '#007bff', textDecoration: 'underline' }}
      >
        {isRegister ? 'Sudah punya akun? Login' : 'Belum punya akun? Daftar'}
      </p>
    </div>
  );
};

// ---------------------------------------------------------
// CREATE TASK FORM (UPDATED WITH RECURRENCE)
// ---------------------------------------------------------

const CreateTaskForm = ({ onTaskCreated }: { onTaskCreated: () => void }) => {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [priority, setPriority] = useState<'low'|'medium'|'high'>('medium');
  const [reminder, setReminder] = useState('');

  // NEW: recurrence
  const [recurrence, setRecurrence] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title) return;

    try {
      const formattedDate = reminder ? new Date(reminder).toISOString() : null;

      await api.post('/tasks/', {
        title,
        description: desc,
        priority,
        reminder_time: formattedDate,
        recurrence_pattern: recurrence   // <---- FIELD BARU
      });

      setTitle('');
      setDesc('');
      setPriority('medium');
      setReminder('');
      setRecurrence('');

      onTaskCreated();
    } catch {
      alert("Gagal membuat task");
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ background: '#f8f9fa', padding: '20px', borderRadius: '12px', marginBottom: '30px' }}>
      <input
        type="text"
        placeholder="Judul Tugas..."
        value={title}
        onChange={e => setTitle(e.target.value)}
        style={{ width: '100%', padding: '12px', borderRadius: '8px', border: '1px solid #ccc', marginBottom: '10px' }}
      />

      <textarea
        placeholder="Deskripsi..."
        value={desc}
        onChange={e => setDesc(e.target.value)}
        style={{ width: '100%', padding: '12px', borderRadius: '8px', border: '1px solid #ccc', minHeight: '60px' }}
      />

      <div style={{ display: 'flex', gap: '10px', marginTop: '10px', flexWrap: 'wrap' }}>
        <select value={priority} onChange={e => setPriority(e.target.value as any)} style={{ padding: '10px', borderRadius: '6px' }}>
          <option value="low">Low Priority</option>
          <option value="medium">Medium Priority</option>
          <option value="high">High Priority</option>
        </select>

        <input
          type="datetime-local"
          value={reminder}
          onChange={e => setReminder(e.target.value)}
          style={{ padding: '10px', borderRadius: '6px' }}
        />

        {/* NEW: recurrence selector */}
        <select
          value={recurrence}
          onChange={e => setRecurrence(e.target.value)}
          style={{ padding: '10px', borderRadius: '6px', border: '1px solid #ccc' }}
        >
          <option value="">Tidak Diulang</option>
          <option value="daily">Setiap Hari</option>
          <option value="weekly">Setiap Minggu</option>
        </select>

        <button type="submit" style={{ padding: '10px 20px', background: '#28a745', color: 'white', borderRadius: '6px', marginLeft: 'auto' }}>
          + Add Task
        </button>
      </div>
    </form>
  );
};

// ---------------------------------------------------------
// SINGLE TASK ITEM (UPDATED WITH 🔄 RECURRING ICON)
// ---------------------------------------------------------

const TaskItem = ({ task, onUpdate, onDelete }: { task: Task, onUpdate: () => void, onDelete: (id: number) => void }) => {
  const priorityColor = { low: '#00C851', medium: '#ffbb33', high: '#ff4444' };

  const calculateProgress = (subs: Subtask[]) => {
    if (!subs || subs.length === 0) return 0;
    return Math.round((subs.filter(s => s.is_done).length / subs.length) * 100);
  };

  const toggleTaskStatus = async () => {
    const newStatus = task.status === 'done' ? 'pending' : 'done';
    await api.put(`/tasks/${task.id}`, { status: newStatus });
    onUpdate();
  };

  const addSubtask = async () => {
    const title = prompt("Nama Checklist:");
    if (!title) return;
    await api.post(`/tasks/${task.id}/subtasks`, { title });
    onUpdate();
  };

  const toggleSubtask = async (id: number) => {
    await api.put(`/subtasks/${id}`);
    onUpdate();
  };

  const deleteSubtask = async (id: number) => {
    if (confirm("Hapus checklist ini?")) {
      await api.delete(`/subtasks/${id}`);
      onUpdate();
    }
  };

  return (
    <div style={{
      background: 'white',
      padding: '20px',
      borderRadius: '12px',
      marginBottom: '15px',
      borderLeft: `6px solid ${priorityColor[task.priority]}`,
      boxShadow: '0 2px 8px rgba(0,0,0,0.05)'
    }}>
      
      {/* HEADER */}
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '10px' }}>
        <div>
          <h3 style={{ margin: 0 }}>
            {/* NEW: Recurring Icon */}
            {task.recurrence_pattern && (
              <span title={`Diulang: ${task.recurrence_pattern}`} style={{ marginRight: '5px' }}>
                🔄
              </span>
            )}

            {task.title}
          </h3>

          <span style={{ fontSize: '0.75em', padding: '4px 8px', borderRadius: '20px', background: '#eee', fontWeight: 'bold' }}>
            {task.priority.toUpperCase()}
          </span>
        </div>

        <div style={{ display: 'flex', gap: '8px' }}>
          <button onClick={toggleTaskStatus} style={{ padding: '6px 12px', borderRadius: '6px' }}>
            {task.status === 'done' ? 'Undo' : 'Done'}
          </button>

          <button onClick={() => onDelete(task.id)} style={{ padding: '6px 12px', background: '#ffebee', color: '#c62828', borderRadius: '6px' }}>
            Del
          </button>
        </div>
      </div>

      {/* DESCRIPTION */}
      {task.description && (
        <p style={{ marginTop: '5px', color: '#666' }}>{task.description}</p>
      )}

      {/* REMINDER */}
      {task.reminder_time && (
        <div style={{ fontSize: '0.85em', color: '#007bff', marginTop: '5px' }}>
          ⏰ {new Date(task.reminder_time).toLocaleString()}
        </div>
      )}

      {/* PROGRESS BAR */}
      {task.subtasks && task.subtasks.length > 0 && (
        <div style={{ marginTop: '10px', background: '#eee', borderRadius: '4px', height: '8px', overflow: 'hidden' }}>
          <div style={{
            height: '100%',
            background: '#28a745',
            width: `${calculateProgress(task.subtasks)}%`,
            transition: 'width .3s'
          }} />
        </div>
      )}

      {/* SUBTASK LIST */}
      <div style={{ marginTop: '10px' }}>
        {task.subtasks?.map(sub => (
          <div key={sub.id} style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '6px' }}>
            <input type="checkbox" checked={sub.is_done} onChange={() => toggleSubtask(sub.id)} />

            <span style={{ textDecoration: sub.is_done ? 'line-through' : 'none', flex: 1 }}>
              {sub.title}
            </span>

            <span onClick={() => deleteSubtask(sub.id)} style={{ color: 'red', cursor: 'pointer' }}>
              ×
            </span>
          </div>
        ))}

        <button
          onClick={addSubtask}
          style={{ background: 'none', border: 'none', color: '#007bff', cursor: 'pointer', marginTop: '5px' }}
        >
          + Add Checklist Item
        </button>
      </div>
    </div>
  );
};

// ---------------------------------------------------------
// MAIN APP
// ---------------------------------------------------------

const ALARM_SOUND_URL = "https://actions.google.com/sounds/v1/alarms/digital_watch_alarm_long.ogg";

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [tasks, setTasks] = useState<Task[]>([]);
  const audioRef = useRef<HTMLAudioElement>(new Audio(ALARM_SOUND_URL));

  // Fetch tasks after login
  useEffect(() => {
    if (token) fetchTasks();
  }, [token]);

  // Alarm checker
  useEffect(() => {
    if (!token || tasks.length === 0) return;

    const interval = setInterval(() => {
      const now = new Date();

      tasks.forEach(task => {
        if (!task.reminder_time || task.status === 'done') return;

        const reminderTime = new Date(task.reminder_time);
        const diff = now.getTime() - reminderTime.getTime();

        if (diff >= 0 && diff < 5000) {
          audioRef.current.play().catch(() => {});
          if (confirm(`⏰ ALARM: "${task.title}"!`)) {
            audioRef.current.pause();
            audioRef.current.currentTime = 0;
          }
        }
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [tasks, token]);

  // Fetch tasks
  const fetchTasks = async () => {
    try {
      const res = await api.get<Task[]>('/tasks/');
      setTasks(res.data || []);
    } catch (err) {
      if ((err as AxiosError).response?.status === 401) handleLogout();
    }
  };

  const handleDeleteTask = async (id: number) => {
    if (confirm("Hapus tugas ini?")) {
      await api.delete(`/tasks/${id}`);
      fetchTasks();
    }
  };

  const handleLoginSuccess = (newToken: string) => {
    localStorage.setItem('token', newToken);
    setToken(newToken);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setTasks([]);
  };

  // Render
  if (!token) return <AuthForm onLogin={handleLoginSuccess} />;

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '30px' }}>
        <h2>Task Dashboard</h2>
        <button onClick={handleLogout} style={{ background: '#dc3545', color: 'white', padding: '8px 16px', borderRadius: '6px' }}>
          Logout
        </button>
      </header>

      <CreateTaskForm onTaskCreated={fetchTasks} />

      <div>
        {tasks.map(task => (
          <TaskItem key={task.id} task={task} onUpdate={fetchTasks} onDelete={handleDeleteTask} />
        ))}

        {tasks.length === 0 && (
          <div style={{ textAlign: 'center', padding: '50px', color: '#888' }}>
            Belum ada tugas. Yuk mulai hari ini! 🚀
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
