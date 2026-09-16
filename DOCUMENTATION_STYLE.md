# 📚 Guía de Estilo de Documentación

Este documento define el estándar de formato para toda la documentación del proyecto **Klaus-antipatterns-search** a partir de septiembre 2026.

## 🎯 Estructura de README.md

### Portada
```markdown
# 🔍 Nombre del Proyecto

> **Tagline en negrita.** Contexto adicional.

[![Badge1](url)](url) [![Badge2](url)](url)

---
```

### Sección central: "¿Qué hago? ¿Cómo lo hago? ¿Y para qué lo hago?"
```markdown
## 🤔 ¿Qué hago? ¿Cómo lo hago? ¿Y para qué lo hago?

**¿Qué hago:** [Párrafo 1 - qué es]
**¿Cómo lo hago:** [Párrafo 2 - técnica]
**¿Para qué lo hago:** [Párrafo 3 - valor/problema]

---
```

### Arquitectura (Flowchart Mermaid preferible)
```markdown
## 🏗️ Arquitectura

\`\`\`mermaid
graph TD
    A["🖥️ Componente A"] --> B["🔌 Componente B"]
    B --> C["📊 Salida"]
\`\`\`

**Stack:**
- Componente 1: tecnología
- Componente 2: tecnología

---
```

### Quick Start / Instalación
```markdown
## 🚀 Instalación / Quick Start

### Requisitos
- [Req 1]

### Pasos
[Código bash con comentarios]

---
```

### Características / Modos
```markdown
## 🎯 Características principales

| Característica | Descripción |
|---|---|
| 🔌 Feature 1 | Descripción breve |
| 💬 Feature 2 | Descripción breve |

---
```

### Secciones adicionales
```markdown
## ⚙️ Configuración

| Variable | Defecto | Descripción |
|---|---|---|
| VAR_1 | valor | qué hace |

---

## 📚 Documentación

| Doc | Contenido |
|---|---|
| [🏗️ File](docs/file.md) | Descripción |

---

## 📊 Roadmap

| Fase | Contenido | Estado |
|---|---|---|
| **v1.0** | Feature A | ✅ Completado |
| **v1.1** | Feature B | 📋 Planeado |

---

## 🤝 Contribuir

1. Fork
2. Rama: `git checkout -b GH-{N}-desc`
3. Commit: `git commit -m "GH-{N}: desc"`
4. Push y PR

---

## 📄 Licencia

[MIT](LICENSE) — Libre para uso comercial.
```

## 🎨 Emojis estándar

| Sección | Emoji |
|---|---|
| Título | 🔍 (proyecto-específico) |
| ¿Qué hago? | 🤔 |
| Arquitectura | 🏗️ |
| Instalación | 🚀 |
| Características | 🎯 |
| Configuración | ⚙️ |
| Documentación | 📚 |
| Roadmap | 🗺️ |
| Tests | 🧪 |
| Contribuir | 🤝 |
| Licencia | 📄 |

## ✅ Checklist de calidad

- [ ] Portada con emoji + tagline + badges (máx 4)
- [ ] Sección "¿Qué hago? ¿Cómo? ¿Para qué?" (3 párrafos)
- [ ] Separadores `---` entre secciones principales
- [ ] Arquitectura con Flowchart Mermaid (si ≥4 componentes)
- [ ] Quick Start con pasos claros y requisitos
- [ ] Características/Features en tabla con emojis distintos
- [ ] Documentación enlazada a `docs/` con emojis
- [ ] Contribuir con pasos numerados y claros
- [ ] Licencia con enlace (preferibly MIT)
- [ ] Emojis consistentes y contextuales
- [ ] Sin párrafos largos (máx 4 líneas cada uno)
- [ ] Código en bloques ```bash``` con comentarios

## 📝 Referencia

Consulta el README.md principal del proyecto como ejemplo de aplicación de este estándar.
