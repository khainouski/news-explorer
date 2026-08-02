BEGIN;

-- The only account that can log in - no public signup (see migration/postgres/20260507170000_schema.up.sql).
-- This is a deliberately simplified/temporary setup: login 'admin', password '12345' (bcrypt cost
-- 10), no email - nothing in the app forces this to change, it's on whoever runs this migration
-- to go change it from /account right after logging in.
-- Inserted first and explicitly as id 1 (internal/domain.AdminUserID) since the app treats "user
-- id 1" as the admin; setval keeps the SERIAL sequence in sync so nothing else ever collides
-- with it.
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

-- Only sources whose feed URL was actually verified reachable (see the articles insert below,
-- which is real content fetched from these feeds) are seeded. Alex Edwards Blog, PostgreSQL News,
-- Reddit r/golang, Remote OK (both), Remote.co Developer Jobs, RemoteYeah Engineering and Indeed
-- Go Jobs all 403/404/timed out when fetched and were dropped rather than kept with a source
-- nothing can ever sync from.
INSERT INTO sources (id, name, feed_url, description, tag_id, badge, badge_color, status, last_synced_at) VALUES
    ('go-blog', 'Go Blog', 'https://go.dev/blog/feed.atom', 'Official Go language announcements, releases, and technical articles.', 'go', 'Go', 'bg-blue-500', 'active', NOW() - INTERVAL '2 minutes'),
    ('golang-weekly', 'Golang Weekly', 'https://golangweekly.com/rss', 'Weekly roundup of Go news, libraries, tutorials, and jobs.', 'go', 'GW', 'bg-sky-600', 'active', NOW() - INTERVAL '6 minutes'),
    ('ardan-labs-blog', 'Ardan Labs Blog', 'https://www.ardanlabs.com/blog/index.xml', 'Advanced Go articles focused on concurrency, performance, and production systems.', 'go', 'AL', 'bg-slate-700', 'active', NOW() - INTERVAL '22 minutes'),
    ('three-dots-labs', 'Three Dots Labs', 'https://threedots.tech/index.xml', 'Go architecture, DDD, CQRS, and microservices.', 'go', '3D', 'bg-purple-600', 'active', NOW() - INTERVAL '40 minutes'),
    ('bitfield-consulting', 'Bitfield Consulting', 'https://bitfieldconsulting.com/posts?format=rss', 'Go language deep dives and best practices.', 'go', 'BC', 'bg-amber-600', 'active', NOW() - INTERVAL '3 hours'),
    ('jetbrains-goland-blog', 'JetBrains GoLand Blog', 'https://blog.jetbrains.com/go/feed/', 'GoLand IDE updates and Go development tips.', 'go', 'GL', 'bg-fuchsia-600', 'active', NOW() - INTERVAL '45 minutes'),
    ('go-time-podcast', 'Go Time Podcast', 'https://changelog.com/gotime/feed', 'Weekly Go podcast covering language news and ecosystem discussions.', 'go', 'GT', 'bg-pink-600', 'active', NOW() - INTERVAL '1 day'),
    ('dev-to-go', 'Dev.to (Go)', 'https://dev.to/feed/tag/golang', 'Community-written Go articles and tutorials.', 'go', '{}', 'bg-gray-900', 'active', NOW() - INTERVAL '8 minutes'),

    ('golang-weekly-jobs', 'Golang Weekly Jobs', 'https://golangweekly.com/rss', 'Weekly Go job opportunities included with the newsletter.', 'jobs', 'GJ', 'bg-sky-700', 'active', NOW() - INTERVAL '6 minutes'),
    ('we-work-remotely', 'We Work Remotely', 'https://weworkremotely.com/categories/remote-programming-jobs.rss', 'Remote software engineering and backend jobs.', 'jobs', 'WWR', 'bg-green-600', 'active', NOW() - INTERVAL '18 minutes'),
    ('hacker-news-jobs', 'Hacker News Jobs', 'https://hnrss.org/jobs', 'Hacker News "Who Is Hiring" and startup job postings.', 'jobs', 'Y', 'bg-orange-500', 'active', NOW() - INTERVAL '5 minutes'),

    ('kubernetes-blog', 'Kubernetes Blog', 'https://kubernetes.io/feed.xml', 'Official Kubernetes releases, features, and best practices.', 'kubernetes', '⎈', 'bg-blue-600', 'active', NOW() - INTERVAL '12 minutes'),
    ('cncf-blog', 'CNCF Blog', 'https://www.cncf.io/feed/', 'Cloud Native ecosystem news including Kubernetes, Prometheus, Argo CD, and OpenTelemetry.', 'cloud-native', 'CN', 'bg-indigo-700', 'active', NOW() - INTERVAL '15 minutes'),

    ('cloudflare-blog', 'Cloudflare Blog', 'https://blog.cloudflare.com/rss/', 'Networking, performance, security, and infrastructure engineering.', 'cloud', 'CF', 'bg-orange-500', 'active', NOW() - INTERVAL '9 minutes'),
    ('aws-news', 'AWS News Blog', 'https://aws.amazon.com/blogs/aws/feed/', 'Official AWS announcements and cloud service updates.', 'cloud', 'AWS', 'bg-gray-900', 'active', NOW() - INTERVAL '15 minutes');

-- Every row here is real content fetched from the source's live feed as of 2026-08-07 (no
-- synthetic placeholders). There is no RSS import job yet (internal/adapter/postgres/article
-- just reads this table) - this is a one-time snapshot, not something that stays fresh on its own.
INSERT INTO articles (id, source_id, external_id, title, summary, url, published_at, unread) VALUES
    ('pkg-go-dev-api', 'go-blog', 'https://go.dev/blog/pkgsite-api', 'Introducing the pkg.go.dev API', 'Introducing the new programmatic API for pkg.go.dev, allowing developers to fetch package and module data directly.', 'https://go.dev/blog/pkgsite-api', '2026-05-21 00:00:00+00', FALSE),
    ('type-construction-cycle-detection', 'go-blog', 'https://go.dev/blog/type-construction-and-cycle-detection', 'Type Construction and Cycle Detection', 'Go 1.26 simplifies type construction and enhances cycle detection for certain kinds of recursive types.', 'https://go.dev/blog/type-construction-and-cycle-detection', '2026-03-24 00:00:00+00', FALSE),
    ('go-fix-inline-source-level-inliner', 'go-blog', 'https://go.dev/blog/inliner', '//go:fix inline and the source-level inliner', 'How Go 1.26''s source-level inliner works, and how it can help you with self-service API migrations.', 'https://go.dev/blog/inliner', '2026-03-10 00:00:00+00', FALSE),

    ('golang-weekly-613', 'golang-weekly', 'https://golangweekly.com/issues/613', 'An interactive tour of Go 1.27', 'An interactive guide covering Go 1.27''s new features, including generic methods, the goroutineleak profile, the uuid package, and experimental portable SIMD support.', 'https://golangweekly.com/issues/613', '2026-08-07 00:00:00+00', TRUE),
    ('golang-weekly-612', 'golang-weekly', 'https://golangweekly.com/issues/612', 'A plan to bring generic collections to Go 1.28', 'The Go Collections working group proposes standard library collection types for Go 1.28: a canonical set type, custom-hasher maps/sets, ordered maps, and a generic binary heap.', 'https://golangweekly.com/issues/612', '2026-07-31 00:00:00+00', FALSE),

    ('ardan-labs-rag-go', 'ardan-labs-blog', 'https://www.ardanlabs.com/blog/2026/04/rag-in-go-a-vulnerability-research-tool/', 'RAG in Go: A Vulnerability Research Tool', 'Retrieval-Augmented Generation lets large language models reach beyond their training data by storing documents with embeddings in a custom database.', 'https://www.ardanlabs.com/blog/2026/04/rag-in-go-a-vulnerability-research-tool/', '2026-04-20 00:00:00+00', FALSE),
    ('ardan-labs-meeting-scheduler', 'ardan-labs-blog', 'https://www.ardanlabs.com/blog/2026/03/using-tools-a-meeting-scheduler/', 'Using Tools: A Meeting Scheduler', 'A walkthrough of function calling - letting LLMs invoke external tools when doing so produces better answers to user queries.', 'https://www.ardanlabs.com/blog/2026/03/using-tools-a-meeting-scheduler/', '2026-03-16 00:00:00+00', FALSE),

    ('threedots-reading-bottleneck', 'three-dots-labs', 'https://threedots.tech/post/understanding-code-is-bottleneck/', 'Writing code isn''t the bottleneck anymore, reading is', 'With more code being AI-generated, code review - not production - has become the primary bottleneck in software development.', 'https://threedots.tech/post/understanding-code-is-bottleneck/', '2026-08-06 00:00:00+02', TRUE),
    ('threedots-distributed-monolith-trap', 'three-dots-labs', 'https://threedots.tech/episode/the-distributed-monolith-trap/', 'The Distributed Monolith Trap (And How to Escape It)', 'Why splitting a monolith into microservices too early creates tight coupling over HTTP - and why starting with a modular monolith usually works better.', 'https://threedots.tech/episode/the-distributed-monolith-trap/', '2026-01-28 16:00:00+00', FALSE),

    ('bitfield-code-like-pirate', 'bitfield-consulting', 'https://bitfieldconsulting.com/posts/code-like-pirate', 'Code like a pirate with AI', 'A six-step PIRATE workflow for using AI coding agents effectively: careful planning, iteration, review, assessment, testing, and evaluation.', 'https://bitfieldconsulting.com/posts/code-like-pirate', '2026-08-04 00:00:00+00', TRUE),
    ('bitfield-castle-building-marketing', 'bitfield-consulting', 'https://bitfieldconsulting.com/posts/castle-building', 'Castle building for beginners: marketing yourself', 'How engineers can market their own products and services by building a website and owned platforms instead of relying on social media.', 'https://bitfieldconsulting.com/posts/castle-building', '2026-07-13 00:00:00+00', FALSE),

    ('jetbrains-go-127-release-party', 'jetbrains-goland-blog', 'https://blog.jetbrains.com/go/2026/08/05/new-livestream-go-127/', 'Go 1.27 Release Party - Free Online Event With the Go Team', 'JetBrains is hosting a free online event on August 25 with Go team members discussing the Go 1.27 release and new language features.', 'https://blog.jetbrains.com/go/2026/08/05/new-livestream-go-127/', '2026-08-05 00:00:00+00', TRUE),
    ('jetbrains-escape-analysis', 'jetbrains-goland-blog', 'https://blog.jetbrains.com/go/2026/07/20/escape-analysis/', 'Escape Analysis in Go - Stack vs. Heap Allocations Explained', 'How Go''s compiler decides whether values are allocated on the stack or the heap, and how GoLand visualizes those decisions.', 'https://blog.jetbrains.com/go/2026/07/20/escape-analysis/', '2026-07-20 00:00:00+00', FALSE),

    ('go-time-340', 'go-time-podcast', 'https://changelog.com/gotime/340', 'That''s Go Time!', 'The finale episode after eight years and 340 episodes, with the full cast reminiscing about favorite moments.', 'https://changelog.com/gotime/340', '2024-12-18 00:00:00+00', FALSE),
    ('go-time-339', 'go-time-podcast', 'https://changelog.com/gotime/339', 'Pitching Go in 2025', 'Exploring when Go remains relevant among newer languages, and how to advocate for its adoption.', 'https://changelog.com/gotime/339', '2024-12-10 00:00:00+00', FALSE),

    ('devto-single-for-loop', 'dev-to-go', 'https://dev.to/shroukabozeid/learning-go-as-a-ruby-developer-part-4-what-can-we-do-with-a-single-for-3320', 'Learning Go as a Ruby Developer - Part 4: What can we do with a single "for"', 'Go consolidates multiple loop constructs into a single `for` keyword, handling traditional iteration, conditional looping, infinite loops, and collection traversal via `range`.', 'https://dev.to/shroukabozeid/learning-go-as-a-ruby-developer-part-4-what-can-we-do-with-a-single-for-3320', '2026-08-07 16:32:53+00', TRUE),
    ('devto-go-plugins-interpreter', 'dev-to-go', 'https://dev.to/alexispires/how-i-made-go-feel-interpreted-with-go-plugins-3pog', 'How I made Go feel interpreted with Go plugins', 'An interactive notebook kernel for Go built on the `plugin` package, dynamically loading compiled cells while maintaining state across executions.', 'https://dev.to/alexispires/how-i-made-go-feel-interpreted-with-go-plugins-3pog', '2026-08-07 14:35:26+00', TRUE),

    ('wwr-gusto-retirement-specialist', 'we-work-remotely', 'https://weworkremotely.com/remote-jobs/gusto-inc-retirement-implementation-specialist', 'Gusto, Inc.: Retirement Implementation Specialist', 'Remote role managing 401(k) plan onboarding for new implementations and conversions.', 'https://weworkremotely.com/remote-jobs/gusto-inc-retirement-implementation-specialist', NOW() - INTERVAL '2 hours', TRUE),
    ('wwr-dropbox-director-product-design', 'we-work-remotely', 'https://weworkremotely.com/remote-jobs/dropbox-director-product-design', 'Dropbox: Director, Product Design', 'Remote leadership role overseeing design systems expansion across Dropbox''s product line.', 'https://weworkremotely.com/remote-jobs/dropbox-director-product-design', NOW() - INTERVAL '5 hours', FALSE),

    ('hn-jobs-truemetrics-gtm-lead', 'hacker-news-jobs', 'https://www.ycombinator.com/companies/truemetrics/jobs/bIQQ7tP-founding-gtm-lead', 'Truemetrics (YC S23) Is Hiring in Berlin - GTM Lead', 'Truemetrics is looking for a founding GTM lead to join their Berlin team.', 'https://www.ycombinator.com/companies/truemetrics/jobs/bIQQ7tP-founding-gtm-lead', '2026-08-04 17:00:00+00', TRUE),
    ('hn-jobs-roame-lead-engineer', 'hacker-news-jobs', 'https://www.ycombinator.com/companies/roame/jobs/mqqfa38-lead-full-stack-engineer', 'Roame (YC S23) Is Hiring Lead Engineer', 'Roame is hiring a lead full-stack engineer to help build out their core product.', 'https://www.ycombinator.com/companies/roame/jobs/mqqfa38-lead-full-stack-engineer', '2026-08-04 12:00:00+00', FALSE),

    ('gateway-api-v1-6-release', 'kubernetes-blog', 'https://kubernetes.io/blog/2026/08/03/gateway-api-v1-6-release/', 'Gateway API v1.6: TCPRoute and UDPRoute Graduate to Standard', 'Gateway API v1.6.0 brings TCPRoute and UDPRoute to GA stability, enabling portable L4 routing for TCP and UDP workloads across Kubernetes.', 'https://kubernetes.io/blog/2026/08/03/gateway-api-v1-6-release/', '2026-08-03 08:00:00-08', TRUE),
    ('kubernetes-v1-37-sneak-peek', 'kubernetes-blog', 'https://kubernetes.io/blog/2026/07/31/kubernetes-v1-37-sneak-peek/', 'Kubernetes v1.37 Sneak Peek', 'A look at what''s landing in Kubernetes v1.37, including the Metrics API graduating to GA and kubelet rootless mode reaching Beta.', 'https://kubernetes.io/blog/2026/07/31/kubernetes-v1-37-sneak-peek/', '2026-07-31 08:00:00-08', TRUE),
    ('controller-runtime-cache-explained', 'kubernetes-blog', 'https://kubernetes.io/blog/2026/07/29/controller-runtime-cache-explained/', 'How the controller-runtime Cache Actually Works', 'A deep dive into controller-runtime''s list-watch cache architecture, and the common mistakes that trip up custom Kubernetes controllers.', 'https://kubernetes.io/blog/2026/07/29/controller-runtime-cache-explained/', '2026-07-29 10:00:00-08', FALSE),

    ('cncf-opencost-inference-tracking', 'cncf-blog', 'https://www.cncf.io/blog/2026/08/05/opencost-1-121-0-first-of-a-kind-kubernetes-inference-cost-tracking/', 'OpenCost 1.121.0: First-of-a-kind Kubernetes inference cost tracking', 'OpenCost now integrates with llm-d to track AI model inference costs, measuring both allocation- and usage-based expenses per model and per token on Kubernetes.', 'https://www.cncf.io/blog/2026/08/05/opencost-1-121-0-first-of-a-kind-kubernetes-inference-cost-tracking/', '2026-08-05 11:00:00+00', TRUE),
    ('cncf-k8gb-incubating', 'cncf-blog', 'https://www.cncf.io/announcements/2026/08/05/k8gb-becomes-a-cncf-incubating-project/', 'K8gb becomes a CNCF incubating project', 'The CNCF Technical Oversight Committee approved Kubernetes Global Balancer (K8gb) as an incubating project for cloud-native global server load balancing.', 'https://www.cncf.io/announcements/2026/08/05/k8gb-becomes-a-cncf-incubating-project/', '2026-08-05 16:00:00+00', TRUE),

    ('cloudflare-agentic-internet', 'cloudflare-blog', 'https://blog.cloudflare.com/good-and-bad-agentic-behaviors/', 'Unveiling good and bad behaviors on the Agentic Internet', 'Cloudflare moves from point-in-time risk assessment to continuous trust evaluation for bot detection, introducing behavioral analysis tools like Precursor.', 'https://blog.cloudflare.com/good-and-bad-agentic-behaviors/', '2026-08-07 13:01:00+00', TRUE),
    ('cloudflare-radar-researcher', 'cloudflare-blog', 'https://blog.cloudflare.com/introducing-radar-researcher/', 'Introducing Radar Researcher: An AI tool for exploring Internet data in plain language', 'A new AI-powered tool lets users query global Internet trends using natural language, generating interactive charts without needing API knowledge.', 'https://blog.cloudflare.com/introducing-radar-researcher/', '2026-08-07 13:00:00+00', FALSE),

    ('aws-bedrock-agentcore-runtime', 'aws-news', 'https://aws.amazon.com/blogs/aws/runtime-instances-persistent-compute-for-production-ai-agents-on-amazon-bedrock-agentcore/', 'Runtime instances: persistent compute for production AI agents on Amazon Bedrock AgentCore', 'AWS introduces persistent, managed EC2 infrastructure for production AI agents on Bedrock AgentCore, supporting multi-agent collaboration and GPU access.', 'https://aws.amazon.com/blogs/aws/runtime-instances-persistent-compute-for-production-ai-agents-on-amazon-bedrock-agentcore/', '2026-08-06 00:00:00+00', TRUE),
    ('aws-dynamodb-vector-search', 'aws-news', 'https://aws.amazon.com/blogs/aws/amazon-dynamodb-now-supports-real-time-vector-search-at-any-scale/', 'Amazon DynamoDB now supports real-time vector search at any scale', 'DynamoDB now provides native vector search with single-digit millisecond latency at 99%+ recall, without a separate vector database.', 'https://aws.amazon.com/blogs/aws/amazon-dynamodb-now-supports-real-time-vector-search-at-any-scale/', '2026-08-05 00:00:00+00', FALSE);

COMMIT;
