BEGIN;

-- admin/12345 - temporary, change it from /account.
INSERT INTO users (id, login, password_hash) VALUES
    (1, 'admin', '$2a$10$2OQNDgAzZWOUa7VVKqffIOAEVP9RnYPn9XCsHXBLwnA5Op3dosKrC');
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

INSERT INTO tags (id, name) VALUES
    ('go', 'Go'),
    ('jobs', 'Jobs'),
    ('kubernetes', 'Kubernetes'),
    ('cloud-native', 'Cloud Native'),
    ('postgresql', 'PostgreSQL'),
    ('cloud', 'Cloud');

INSERT INTO sources (id, name, feed_url, description, tag_id, badge, badge_color, status, last_synced_at) VALUES
    ('go-blog', 'Go Blog', 'https://go.dev/blog/feed.atom', 'Official Go language announcements, releases, and technical articles.', 'go', 'Go', 'bg-blue-500', 'active', NULL),
    ('golang-weekly', 'Golang Weekly', 'https://golangweekly.com/rss', 'Weekly roundup of Go news, libraries, tutorials, and jobs.', 'go', 'GW', 'bg-sky-600', 'active', NULL),
    ('ardan-labs-blog', 'Ardan Labs Blog', 'https://www.ardanlabs.com/blog/index.xml', 'Advanced Go articles focused on concurrency, performance, and production systems.', 'go', 'AL', 'bg-slate-700', 'active', NULL),
    ('three-dots-labs', 'Three Dots Labs', 'https://threedots.tech/index.xml', 'Go architecture, DDD, CQRS, and microservices.', 'go', '3D', 'bg-purple-600', 'active', NULL),
    ('bitfield-consulting', 'Bitfield Consulting', 'https://bitfieldconsulting.com/posts?format=rss', 'Go language deep dives and best practices.', 'go', 'BC', 'bg-amber-600', 'active', NULL),
    ('jetbrains-goland-blog', 'JetBrains GoLand Blog', 'https://blog.jetbrains.com/go/feed/', 'GoLand IDE updates and Go development tips.', 'go', 'GL', 'bg-fuchsia-600', 'active', NULL),
    ('go-time-podcast', 'Go Time Podcast', 'https://changelog.com/gotime/feed', 'Weekly Go podcast covering language news and ecosystem discussions.', 'go', 'GT', 'bg-pink-600', 'active', NULL),
    ('dev-to-go', 'Dev.to (Go)', 'https://dev.to/feed/tag/golang', 'Community-written Go articles and tutorials.', 'go', '{}', 'bg-gray-900', 'active', NULL),

    ('golang-weekly-jobs', 'Golang Weekly Jobs', 'https://golangweekly.com/rss', 'Weekly Go job opportunities included with the newsletter.', 'jobs', 'GJ', 'bg-sky-700', 'active', NULL),
    ('we-work-remotely', 'We Work Remotely', 'https://weworkremotely.com/categories/remote-programming-jobs.rss', 'Remote software engineering and backend jobs.', 'jobs', 'WWR', 'bg-green-600', 'active', NULL),

    ('kubernetes-blog', 'Kubernetes Blog', 'https://kubernetes.io/feed.xml', 'Official Kubernetes releases, features, and best practices.', 'kubernetes', '⎈', 'bg-blue-600', 'active', NULL),
    ('cncf-blog', 'CNCF Blog', 'https://www.cncf.io/feed/', 'Cloud Native ecosystem news including Kubernetes, Prometheus, Argo CD, and OpenTelemetry.', 'cloud-native', 'CN', 'bg-indigo-700', 'active', NULL),

    ('cloudflare-blog', 'Cloudflare Blog', 'https://blog.cloudflare.com/rss/', 'Networking, performance, security, and infrastructure engineering.', 'cloud', 'CF', 'bg-orange-500', 'active', NULL),
    ('aws-news', 'AWS News Blog', 'https://aws.amazon.com/blogs/aws/feed/', 'Official AWS announcements and cloud service updates.', 'cloud', 'AWS', 'bg-gray-900', 'active', NULL);

COMMIT;
