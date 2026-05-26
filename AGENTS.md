# Lifeline Project Guidelines

## Project Overview

Lifeline is an activity tracker built with Go, React, Vite, PostgreSQL, and Docker.

## Architecture

The codebase follows DDD and is split into four layers:

- `domain`: core business entities and domain logic. This layer must not depend on other layers.
- `infrastructure`: technical integrations such as databases, networking, and internal platform concerns.
- `application`: use cases that coordinate domain logic and infrastructure to execute business scenarios.
- `presentation`: entry points for external systems. This layer validates request shape and delegates to `application`.

## Business Model

Key business concepts:

- `habit`: created once by a user and applied to every day.
- `entry`: created once per day by a user and contains mood, completed habits, diary text, and photos.

## Git And Change Management

- Use Conventional Commits, for example: `feat(habits): add habit completion tracking`.
- Do not use imperative-style commit messages.
- Split changes into semantic groups and create a separate commit for each feature or concern.
- Do not commit local IDE files such as `dataSources.xml`.

## Application And Use Case Rules

- Each use case must live in its own Go package.
- A use case handler is a struct created through a constructor and exposing a single `Handle` method.
- Return simple types unless a dedicated result type provides real value.
- Keep error handling and guard clauses before the happy path.
- Wrap errors with method and layer context instead of returning bare `err`.
- `main.go` is an orchestration file only: construct dependencies and wire them together.
- Split `main.go` into semantic sections with comments such as `--- Section: Usecase handlers ---`.

## HTTP Layer Rules

- Keep the HTTP layer thin: accept the request, validate structural integrity, call the use case, map the result to the response.
- Business validation belongs in `domain` or `application`, not in `presentation`.
- Apply middlewares in `main.go`, not inside route handlers.
- Internal HTTP handlers must use the signature `func(request *http.Request) (status int, payload any)`.
- Each HTTP route must have its own `example.http` file named after the route, for example `PUT /api/v1/habits/{uuid}`.

## Persistence And Identifiers

- Use UUIDv7 instead of numeric identifiers.
- Repository methods must use the `nil, nil` contract to represent "entity not found".
- Error messages should be as specific as possible and include entity identifiers when available.
- Each migration file must operate on exactly one database table.
- Do not use `ON DELETE` clauses in migrations.

## Frontend Conventions

- Use BEM for HTML and React class naming.
- Use 4 spaces for indentation in `html`, `css`, `sass`, `ts`, `js`, and `tsx` files.
