import { useState, useEffect, useRef } from 'react';
import api from './api';
import type { Task, LoginResponse} from './types';
import { AxiosError } from 'axios';
import { 
  CheckCircle2, Circle, Trash2, Plus, LogOut, 
  Calendar, Repeat, AlertCircle, ListTodo,
  LayoutDashboard, CheckSquare, Clock
} from 'lucide-react';

// --- COMPONENTS ---

// 1. Auth Form (Modern Glassy Look)
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
        const res = await api.post<LoginResponse>('/login', { email: formData.email, password: formData.password });
        onLogin(res.data.access_token);
      }
    } catch (error) {
      alert("Gagal! Cek email/password.");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800 p-4">
      <div className="w-full max-w-md bg-slate-800/50 backdrop-blur-xl border border-slate-700 p-8 rounded-2xl shadow-2xl">
        <div className="flex justify-center mb-6">
          <div className="p-3 bg-indigo-500/20 rounded-full">
            <LayoutDashboard size={32} className="text-indigo-400" />
          </div>
        </div>
        <h2 className="text-2xl font-bold text-center mb-2 text-white">
          {isRegister ? 'Create Account' : 'Welcome Back'}
        </h2>
        <p className="text-slate-400 text-center mb-8">Manage your tasks efficiently</p>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {isRegister && (
            <input 
              type="text" placeholder="Full Name" required
              value={formData.name} 
              onChange={e => setFormData({ ...formData, name: e.target.value })}
              className="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 transition"
            />
          )}
          <input 
            type="email" placeholder="Email Address" required
            value={formData.email} 
            onChange={e => setFormData({ ...formData, email: e.target.value })}
            className="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 transition"
          />
          <input 
            type="password" placeholder="Password" required
            value={formData.password} 
            onChange={e => setFormData({ ...formData, password: e.target.value })}
            className="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 transition"
          />
          <button type="submit" className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 rounded-lg transition-all shadow-lg shadow-indigo-500/30">
            {isRegister ? 'Sign Up' : 'Sign In'}
          </button>
        </form>
        <p className="mt-6 text-center text-slate-400 text-sm">
          {isRegister ? "Already have an account? " : "Don't have an account? "}
          <button onClick={() => setIsRegister(!isRegister)} className="text-indigo-400 hover:text-indigo-300 font-medium">
            {isRegister ? "Login" : "Register"}
          </button>
        </p>
      </div>
    </div>
  );
};

// 2. Create Task Form
const CreateTaskForm = ({ onTaskCreated }: { onTaskCreated: () => void }) => {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [priority, setPriority] = useState<'low'|'medium'|'high'>('medium');
  const [recurrence, setRecurrence] = useState('');
  const [reminder, setReminder] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title) return;
    try {
      const formattedDate = reminder ? new Date(reminder).toISOString() : null;
      await api.post('/tasks/', { 
        title, description: desc, priority, 
        reminder_time: formattedDate, recurrence_pattern: recurrence 
      });
      setTitle(''); setDesc(''); setPriority('medium'); setReminder(''); setRecurrence('');
      onTaskCreated();
    } catch (e) { alert("Error creating task"); }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-slate-800 border border-slate-700 rounded-xl p-6 mb-8 shadow-lg">
      <div className="mb-4">
        <input 
          type="text" placeholder="What needs to be done?" value={title} onChange={e => setTitle(e.target.value)}
          className="w-full bg-transparent text-xl font-medium text-white placeholder-slate-500 focus:outline-none mb-2"
        />
        <textarea 
          placeholder="Add a description..." value={desc} onChange={e => setDesc(e.target.value)}
          className="w-full bg-transparent text-sm text-slate-400 placeholder-slate-600 focus:outline-none resize-none h-10"
        />
      </div>
      <div className="flex flex-wrap gap-3 items-center pt-4 border-t border-slate-700">
        <div className="flex items-center gap-2 bg-slate-900 rounded-lg px-3 py-1.5 border border-slate-700">
          <AlertCircle size={16} className={priority === 'high' ? 'text-red-500' : priority === 'medium' ? 'text-amber-500' : 'text-green-500'} />
          <select 
            value={priority} onChange={e => setPriority(e.target.value as any)} 
            className="bg-transparent text-sm text-slate-300 focus:outline-none cursor-pointer"
          >
            <option value="low">Low Priority</option>
            <option value="medium">Medium Priority</option>
            <option value="high">High Priority</option>
          </select>
        </div>

        <div className="flex items-center gap-2 bg-slate-900 rounded-lg px-3 py-1.5 border border-slate-700">
          <Repeat size={16} className="text-blue-400" />
          <select 
            value={recurrence} onChange={e => setRecurrence(e.target.value)} 
            className="bg-transparent text-sm text-slate-300 focus:outline-none cursor-pointer"
          >
            <option value="">One-time</option>
            <option value="daily">Daily</option>
            <option value="weekly">Weekly</option>
          </select>
        </div>

        <div className="flex items-center gap-2 bg-slate-900 rounded-lg px-3 py-1.5 border border-slate-700">
          <Calendar size={16} className="text-purple-400" />
          <input 
            type="datetime-local" value={reminder} onChange={e => setReminder(e.target.value)} 
            className="bg-transparent text-sm text-slate-300 focus:outline-none cursor-pointer"
          />
        </div>

        <button type="submit" className="ml-auto bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition flex items-center gap-2">
          <Plus size={18} /> Add Task
        </button>
      </div>
    </form>
  );
};

// 3. Task Item Component
const TaskItem = ({ task, onUpdate, onDelete }: { task: Task, onUpdate: () => void, onDelete: (id: number) => void }) => {
  const priorityColors = {
    low: 'border-l-4 border-green-500',
    medium: 'border-l-4 border-amber-500',
    high: 'border-l-4 border-red-500',
  };

  const progress = task.subtasks?.length 
    ? Math.round((task.subtasks.filter(s => s.is_done).length / task.subtasks.length) * 100) 
    : 0;

  // Handlers
  const handleToggleStatus = async () => {
    const newStatus = task.status === 'done' ? 'pending' : 'done';
    await api.put(`/tasks/${task.id}`, { status: newStatus });
    onUpdate();
  };

  const handleAddSubtask = async () => {
    const title = prompt("New subtask:");
    if (!title) return;
    await api.post(`/tasks/${task.id}/subtasks`, { title });
    onUpdate();
  };

  const handleToggleSubtask = async (id: number) => {
    await api.put(`/subtasks/${id}`);
    onUpdate();
  };

  return (
    <div className={`bg-slate-800 rounded-xl p-5 mb-4 shadow-sm hover:shadow-md transition border border-slate-700 ${priorityColors[task.priority]}`}>
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-1">
            <h3 className={`text-lg font-semibold ${task.status === 'done' ? 'text-slate-500 line-through' : 'text-white'}`}>
              {task.title}
            </h3>
            {task.recurrence_pattern && (
              <span className="bg-blue-500/20 text-blue-300 text-xs px-2 py-0.5 rounded-full flex items-center gap-1">
                <Repeat size={10} /> {task.recurrence_pattern}
              </span>
            )}
            <span className={`text-xs px-2 py-0.5 rounded-full uppercase font-bold tracking-wider 
              ${task.priority === 'high' ? 'bg-red-500/20 text-red-400' : task.priority === 'medium' ? 'bg-amber-500/20 text-amber-400' : 'bg-green-500/20 text-green-400'}`}>
              {task.priority}
            </span>
          </div>
          
          {task.description && <p className="text-slate-400 text-sm mb-3">{task.description}</p>}
          
          {task.reminder_time && (
            <div className="flex items-center gap-2 text-xs text-indigo-300 bg-indigo-900/30 w-fit px-2 py-1 rounded mb-3">
              <Clock size={12} />
              {new Date(task.reminder_time).toLocaleString()}
            </div>
          )}

          {/* Progress Bar */}
          {task.subtasks?.length > 0 && (
            <div className="mb-4">
              <div className="flex justify-between text-xs text-slate-400 mb-1">
                <span>Progress</span>
                <span>{progress}%</span>
              </div>
              <div className="w-full bg-slate-700 rounded-full h-1.5">
                <div className="bg-indigo-500 h-1.5 rounded-full transition-all duration-300" style={{ width: `${progress}%` }}></div>
              </div>
            </div>
          )}

          {/* Subtasks */}
          <div className="space-y-2">
            {task.subtasks?.map(sub => (
              <div key={sub.id} className="flex items-center gap-3 text-sm group">
                <button onClick={() => handleToggleSubtask(sub.id)} className="text-slate-400 hover:text-indigo-400 transition">
                  {sub.is_done ? <CheckSquare size={16} className="text-indigo-500" /> : <div className="w-4 h-4 border-2 border-slate-500 rounded sm" />}
                </button>
                <span className={sub.is_done ? 'text-slate-500 line-through' : 'text-slate-300'}>{sub.title}</span>
                <button onClick={() => { if(confirm('Delete?')) {api.delete(`/subtasks/${sub.id}`).then(onUpdate)} }} 
                  className="opacity-0 group-hover:opacity-100 text-slate-600 hover:text-red-400 transition ml-auto">
                  <Trash2 size={14} />
                </button>
              </div>
            ))}
            <button onClick={handleAddSubtask} className="text-xs text-indigo-400 hover:text-indigo-300 mt-2 font-medium flex items-center gap-1">
              <Plus size={12} /> Add Checklist Item
            </button>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-col gap-2">
          <button onClick={handleToggleStatus} className={`p-2 rounded-lg transition ${task.status === 'done' ? 'bg-indigo-500/20 text-indigo-300' : 'bg-slate-700 text-slate-400 hover:bg-slate-600'}`}>
            {task.status === 'done' ? <CheckCircle2 size={20} /> : <Circle size={20} />}
          </button>
          <button onClick={() => onDelete(task.id)} className="p-2 rounded-lg bg-slate-700 text-slate-400 hover:bg-red-500/20 hover:text-red-400 transition">
            <Trash2 size={20} />
          </button>
        </div>
      </div>
    </div>
  );
};

// 4. Stats Component
const Stats = ({ tasks }: { tasks: Task[] }) => {
  const total = tasks.length;
  const done = tasks.filter(t => t.status === 'done').length;
  const pending = total - done;

  return (
    <div className="grid grid-cols-3 gap-4 mb-8">
      <div className="bg-slate-800 p-4 rounded-xl border border-slate-700 text-center">
        <div className="text-3xl font-bold text-white mb-1">{total}</div>
        <div className="text-xs text-slate-400 uppercase tracking-wider">Total Tasks</div>
      </div>
      <div className="bg-slate-800 p-4 rounded-xl border border-slate-700 text-center">
        <div className="text-3xl font-bold text-amber-400 mb-1">{pending}</div>
        <div className="text-xs text-slate-400 uppercase tracking-wider">Pending</div>
      </div>
      <div className="bg-slate-800 p-4 rounded-xl border border-slate-700 text-center">
        <div className="text-3xl font-bold text-green-400 mb-1">{done}</div>
        <div className="text-xs text-slate-400 uppercase tracking-wider">Completed</div>
      </div>
    </div>
  );
}

// --- APP ---

const ALARM_SOUND_URL = "https://actions.google.com/sounds/v1/alarms/digital_watch_alarm_long.ogg";

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [tasks, setTasks] = useState<Task[]>([]);
  const audioRef = useRef<HTMLAudioElement>(new Audio(ALARM_SOUND_URL));

  useEffect(() => { if (token) fetchTasks(); }, [token]);

  // Alarm Logic
  useEffect(() => {
    if (!token || tasks.length === 0) return;
    const interval = setInterval(() => {
      const now = new Date();
      tasks.forEach(task => {
        if (task.reminder_time && task.status !== 'done') {
          const reminderTime = new Date(task.reminder_time);
          const diff = now.getTime() - reminderTime.getTime();
          if (diff >= 0 && diff < 5000) {
            audioRef.current.play().catch(console.error);
            if (confirm(`⏰ ALARM: ${task.title}`)) {
              audioRef.current.pause();
              audioRef.current.currentTime = 0;
            }
          }
        }
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [tasks, token]);

  const fetchTasks = async () => {
    try {
      const res = await api.get<Task[]>('/tasks/');
      setTasks(res.data || []);
    } catch (err) {
      if ((err as AxiosError).response?.status === 401) handleLogout();
    }
  };

  const handleDeleteTask = async (id: number) => {
    if (confirm("Delete this task?")) {
      await api.delete(`/tasks/${id}`);
      fetchTasks();
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setTasks([]);
  };

  if (!token) return <AuthForm onLogin={(t) => { localStorage.setItem('token', t); setToken(t); }} />;

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 pb-20">
      {/* Navbar */}
      <nav className="bg-slate-800/50 backdrop-blur-md border-b border-slate-700 sticky top-0 z-10">
        <div className="max-w-3xl mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <div className="bg-indigo-600 p-1.5 rounded-lg">
              <ListTodo size={24} className="text-white" />
            </div>
            <h1 className="text-xl font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
              TaskFlow
            </h1>
          </div>
          <button onClick={handleLogout} className="text-slate-400 hover:text-red-400 transition">
            <LogOut size={20} />
          </button>
        </div>
      </nav>

      {/* Content */}
      <main className="max-w-3xl mx-auto px-4 py-8">
        <Stats tasks={tasks} />
        
        <CreateTaskForm onTaskCreated={fetchTasks} />

        <div className="space-y-4">
          <h2 className="text-sm font-bold text-slate-500 uppercase tracking-widest mb-4">Your Tasks</h2>
          {tasks.map(task => (
            <TaskItem 
              key={task.id} task={task} 
              onUpdate={fetchTasks} onDelete={handleDeleteTask} 
            />
          ))}
          
          {tasks.length === 0 && (
            <div className="text-center py-20 opacity-50">
              <ListTodo size={48} className="mx-auto mb-4 text-slate-600" />
              <p>No tasks yet. Start by creating one!</p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;