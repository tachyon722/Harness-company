-- 002_seed_roles.sql

INSERT INTO roles (name, role_type, description, permissions) VALUES
  ('Strategic Planner', 'planner', 'C-suite and strategic decision makers', '["org:read","org:write","strategy:full","governance:full"]'),
  ('Tactical Planner', 'planner', 'MVRU leads and product managers', '["mvru:read","mvru:write","workflow:full","capability:read"]'),
  ('AI Planner', 'planner', 'AI agents responsible for planning', '["mvru:read","workflow:read","capability:read"]'),
  ('Human Executor', 'executor', 'Human team members executing tasks', '["task:read","task:write","capability:use"]'),
  ('AI Executor', 'executor', 'AI agents executing defined tasks', '["task:read","task:write","capability:use"]'),
  ('Independent Reviewer', 'reviewer', 'Independent reviewers (human)', '["review:full","verification:read"]'),
  ('AI Reviewer', 'reviewer', 'AI agents performing automated review', '["review:limited","verification:read"]')
ON CONFLICT (name) DO NOTHING;
