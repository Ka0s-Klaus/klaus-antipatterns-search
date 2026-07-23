# 🔍 Klaus antipatterns-search

> Multi-language anti-pattern detector — CLI, GitHub Action and multi-org scanner powered by Go + tree-sitter

[![K*](https://img.shields.io/badge/K%2A-AI%20Workspace-7057ff)](https://github.com/Ka0s-Klaus)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Kanban](https://img.shields.io/badge/Kanban-%2321-blue)](https://github.com/orgs/Ka0s-Klaus/projects/21)

---

## 🤔 ¿Qué hago? ¿Cómo lo hago? ¿Y para qué lo hago?

### ¿Qué hago?

Soy un detector de **anti-patrones de código** multi-lenguaje. Analizo repositorios — localmente, en CI/CD o a escala de organización — y genero diagnósticos accionables de deuda técnica: God Objects, complejidad excesiva, código duplicado, dependencias circulares, magic numbers y más.

### ¿Cómo lo hago?

Combino dos enfoques complementarios:

1. **Detectores nativos** sobre AST unificado via [`tree-sitter`](https://tree-sitter.github.io/tree-sitter/) — un único motor de análisis para ~40+ lenguajes.
2. **Adaptadores a linters OSS** (`jscpd`, `radon`, `gocyclo`, `madge`, etc.) — orquestados por subproceso con degradación elegante si no están instalados.

Todos los hallazgos se normalizan al mismo modelo `Finding` y se renderizan en múltiples formatos.

### ¿Para qué lo hago?

Porque trabajar con muchos repositorios en organizaciones distintas (Ka0s-Klaus, mango, masorange) sin una herramienta unificada hace que la deuda técnica se acumule silenciosamente. `Klaus-antipatterns-search` da visibilidad reproducible y consistente sobre la salud del código en cualquier stack.

---

## 🏗️ Arquitectura

```mermaid
graph TD
    CLI["🖥️ CLI / GitHub Action"] --> CORE["🎛️ Core Orchestrator"]
    CORE --> LANG["🔍 Language Detector"]
    LANG --> NATIVE["🌳 Native Detectors\n(tree-sitter)"]
    LANG --> ADAPTERS["🔌 OSS Adapters\n(jscpd · radon · gocyclo · madge)"]
    NATIVE --> MODEL["📦 model.Finding"]
    ADAPTERS --> MODEL
    MODEL --> REPORT["📊 Report Renderers"]
    REPORT --> CONSOLE["💻 Console + JSON"]
    REPORT --> MARKDOWN["📄 Markdown / HTML"]
    REPORT --> SARIF["🔐 SARIF\n(GitHub Code Scanning)"]

    SCAN["🌐 Multi-org Scanner"] --> CORE
    CONFIG["⚙️ .antipatterns.yml"] --> CORE
```

---

## 🚀 Modos de ejecución

| Modo | Comando | Descripción |
| --- | --- | --- |
| 📁 **Local** | `antipatterns scan ./ruta` | Analiza un repo local |
| 🔁 **CI/CD** | GitHub Action | Comenta en PR + sube SARIF |
| 🌐 **Multi-org** | `antipatterns scan-org <perfil>` | Escanea toda una organización |

---

## 🎯 Anti-patrones detectados (MVP)

| Anti-patrón | Método | Fuente |
| --- | --- | --- |
| 🐘 God Object / God Class | métricas de clase (LOC, métodos, fan-in/out) | nativo (tree-sitter) |
| 🌀 Complejidad excesiva | complejidad ciclomática / cognitiva | OSS: `radon`, `gocyclo`, `eslint` |
| 🔁 Código duplicado | detección de clones tipo-1/tipo-2 | OSS: `jscpd` |
| 🔄 Dependencias circulares | grafo de imports + detección de ciclos | OSS: `madge`, `go list` |
| 🔢 Magic numbers/strings | literales repetidos fuera de constantes | nativo |
| 📏 Funciones/ficheros gigantes | LOC y nº parámetros sobre umbral | nativo |
| 💀 Código muerto (lava flow) | símbolos sin referencias | OSS: `vulture`, `ts-prune` |
| 📋 Copy-paste programming | agregación de duplicación por directorio | derivado |

---

## 📊 Formatos de reporte

- **Terminal** — tabla coloreada por severidad (info/low/medium/high/critical)
- **JSON** — `report.json` estructurado para integración con otras herramientas
- **Markdown/HTML** — reporte navegable con métricas y enlaces a líneas afectadas
- **SARIF** — integración con GitHub Code Scanning (pestaña Security)

---

## ⚙️ Configuración

```yaml
# .antipatterns.yml
orgs:
  ka0s-klaus:
    token_env: GH_TOKEN_KA0S
    output: reports/ka0s/
    publish: true
  masorange:
    token_env: GH_TOKEN_MASORANGE
    output: reports/masorange/
    publish: false   # datos de cliente: nunca se publican

thresholds:
  god_object:
    methods: 20
    loc: 400
  function_loc: 80
  cyclomatic: 15
  duplication_pct: 5
```

---

## 🗺️ Roadmap

| Fase | Contenido | Estado |
| --- | --- | --- |
| **0 — Andamiaje** | módulo Go, CLI Cobra, modelo Finding, renderers console+JSON | 🔄 En progreso |
| **1 — MVP nativo** | tree-sitter + detectores: God Object, funciones gigantes, magic numbers | ⏳ Pendiente |
| **2 — Adaptadores OSS** | `jscpd`, `madge`, `radon`, `gocyclo`, skip elegante | ⏳ Pendiente |
| **3 — SARIF + Action** | GitHub Action, comentario en PR, Code Scanning | ⏳ Pendiente |
| **4 — Multi-org** | `scan-org`, perfiles, paralelismo, panel agregado | ⏳ Pendiente |
| **5 — Publicación** | README, licencia, releases cross-compilados | ⏳ Pendiente |

---

## 🤝 Contribuir

Este proyecto sigue el flujo **Issue → Rama → PR** sin excepciones.

1. Abre una Issue con el contexto completo
2. Crea una rama `GH-{N}-descripcion-breve`
3. Trabaja en la rama
4. Abre PR contra `main`

---

## 📄 Licencia

MIT © [Ka0s-Klaus](https://github.com/Ka0s-Klaus)
