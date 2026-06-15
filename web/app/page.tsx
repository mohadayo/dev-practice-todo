"use client";

import { useEffect, useState } from "react";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

type Task = {
  id: number;
  title: string;
  done: boolean;
  created_at: string;
};

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState("");

  async function load() {
    const res = await fetch(`${API_BASE}/tasks`);
    setTasks(await res.json());
  }

  useEffect(() => {
    load();
  }, []);

  async function addTask(e: React.FormEvent) {
    e.preventDefault();
    const t = title.trim();
    if (!t) return;
    await fetch(`${API_BASE}/tasks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title: t }),
    });
    setTitle("");
    await load();
  }

  async function toggle(task: Task) {
    await fetch(`${API_BASE}/tasks/${task.id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ done: !task.done }),
    });
    await load();
  }

  async function remove(task: Task) {
    await fetch(`${API_BASE}/tasks/${task.id}`, { method: "DELETE" });
    await load();
  }

  return (
    <main>
      <h1>TODO</h1>

      <form onSubmit={addTask} style={{ display: "flex", gap: 8 }}>
        <input
          aria-label="new task title"
          data-testid="new-task-input"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="やることを入力"
          style={{ flex: 1, padding: 8 }}
        />
        <button type="submit" data-testid="add-button" style={{ padding: "8px 16px" }}>
          追加
        </button>
      </form>

      <ul data-testid="task-list" style={{ listStyle: "none", padding: 0, marginTop: 24 }}>
        {tasks.map((task) => (
          <li
            key={task.id}
            data-testid="task-item"
            style={{ display: "flex", alignItems: "center", gap: 8, padding: "6px 0" }}
          >
            <input
              type="checkbox"
              checked={task.done}
              onChange={() => toggle(task)}
              aria-label={`toggle ${task.title}`}
            />
            <span style={{ flex: 1, textDecoration: task.done ? "line-through" : "none" }}>
              {task.title}
            </span>
            <button onClick={() => remove(task)} aria-label={`delete ${task.title}`}>
              削除
            </button>
          </li>
        ))}
      </ul>
    </main>
  );
}
