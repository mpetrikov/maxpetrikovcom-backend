INSERT INTO labs (
    id,
    slug,
    title,
    description,
    difficulty,
    timeout_minutes,
    image,
    cpu_limit,
    memory_limit,
    is_published
)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'linux-processes',
    'Linux Processes',
    'Demo lab for local Kubernetes runtime testing.',
    'easy',
    60,
    'ubuntu:24.04',
    '500m',
    '256Mi',
    true
)
ON CONFLICT (slug) DO UPDATE
SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    difficulty = EXCLUDED.difficulty,
    timeout_minutes = EXCLUDED.timeout_minutes,
    image = EXCLUDED.image,
    cpu_limit = EXCLUDED.cpu_limit,
    memory_limit = EXCLUDED.memory_limit,
    is_published = EXCLUDED.is_published,
    updated_at = NOW();
