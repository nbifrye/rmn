import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'rmn',
  description: 'Fast, open-source command-line client for Redmine written in Go. Includes MCP server for AI agent integration.',
  base: '/rmn/',

  head: [
    ['link', { rel: 'icon', href: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>⚡</text></svg>" }],

    // Open Graph
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:url', content: 'https://nbifrye.github.io/rmn/' }],
    ['meta', { property: 'og:title', content: 'rmn — Redmine CLI Tool' }],
    ['meta', { property: 'og:description', content: 'Fast, open-source command-line client for Redmine written in Go. Includes MCP server for AI agent integration.' }],
    ['meta', { property: 'og:site_name', content: 'rmn' }],

    // Twitter Card
    ['meta', { name: 'twitter:card', content: 'summary' }],
    ['meta', { name: 'twitter:title', content: 'rmn — Redmine CLI Tool' }],
    ['meta', { name: 'twitter:description', content: 'Fast, open-source command-line client for Redmine written in Go. Includes MCP server for AI agent integration.' }],

    // AI Agent Discovery
    ['link', { rel: 'alternate', type: 'text/plain', href: 'https://nbifrye.github.io/rmn/llms.txt', title: 'LLM-friendly documentation' }],

    // JSON-LD Structured Data
    ['script', { type: 'application/ld+json' }, JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'SoftwareApplication',
      name: 'rmn',
      description: 'Command-line client for Redmine written in Go with MCP server support for AI agents',
      url: 'https://nbifrye.github.io/rmn/',
      applicationCategory: 'DeveloperApplication',
      operatingSystem: 'Linux, macOS, Windows',
      programmingLanguage: 'Go',
      license: 'https://opensource.org/licenses/MIT',
      codeRepository: 'https://github.com/nbifrye/rmn',
      offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' }
    })],

    // SEO
    ['meta', { name: 'keywords', content: 'redmine, cli, command-line, go, golang, issue tracker, mcp, model context protocol, ai agent, redmine api, redmine client, project management' }],
  ],

  sitemap: {
    hostname: 'https://nbifrye.github.io/rmn/'
  },

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      themeConfig: {
        nav: [
          { text: 'Guide', link: '/guide/installation' },
          { text: 'MCP Server', link: '/mcp-server' },
          { text: 'Reference', link: '/reference/shell-completion' },
          {
            text: 'Links',
            items: [
              { text: 'GitHub', link: 'https://github.com/nbifrye/rmn' },
              { text: 'Releases', link: 'https://github.com/nbifrye/rmn/releases' },
              { text: 'llms.txt', link: 'https://nbifrye.github.io/rmn/llms.txt' },
            ]
          }
        ],

        sidebar: [
          {
            text: 'Guide',
            items: [
              { text: 'Installation', link: '/guide/installation' },
              { text: 'Configuration', link: '/guide/configuration' },
              { text: 'Usage', link: '/guide/usage' },
            ]
          },
          {
            text: 'Integrations',
            items: [
              { text: 'MCP Server', link: '/mcp-server' },
            ]
          },
          {
            text: 'Reference',
            items: [
              { text: 'Shell Completion', link: '/reference/shell-completion' },
              { text: 'Security', link: '/reference/security' },
              { text: 'Architecture', link: '/reference/architecture' },
            ]
          },
          {
            text: 'Community',
            items: [
              { text: 'Development', link: '/development' },
            ]
          }
        ],

        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © nbifrye'
        },

        editLink: {
          pattern: 'https://github.com/nbifrye/rmn/edit/main/docs/:path',
          text: 'Edit this page on GitHub'
        }
      }
    },
    ja: {
      label: '日本語',
      lang: 'ja',
      title: 'rmn',
      description: 'Go言語で書かれた高速なオープンソースのRedmineコマンドラインクライアント。AIエージェント連携用MCPサーバー搭載。',
      themeConfig: {
        nav: [
          { text: 'ガイド', link: '/ja/guide/installation' },
          { text: 'MCPサーバー', link: '/ja/mcp-server' },
          { text: 'リファレンス', link: '/ja/reference/shell-completion' },
          {
            text: 'リンク',
            items: [
              { text: 'GitHub', link: 'https://github.com/nbifrye/rmn' },
              { text: 'リリース', link: 'https://github.com/nbifrye/rmn/releases' },
              { text: 'llms.txt', link: 'https://nbifrye.github.io/rmn/llms-ja.txt' },
            ]
          }
        ],

        sidebar: [
          {
            text: 'ガイド',
            items: [
              { text: 'インストール', link: '/ja/guide/installation' },
              { text: '設定', link: '/ja/guide/configuration' },
              { text: '使い方', link: '/ja/guide/usage' },
            ]
          },
          {
            text: '連携',
            items: [
              { text: 'MCPサーバー', link: '/ja/mcp-server' },
            ]
          },
          {
            text: 'リファレンス',
            items: [
              { text: 'シェル補完', link: '/ja/reference/shell-completion' },
              { text: 'セキュリティ', link: '/ja/reference/security' },
              { text: 'アーキテクチャ', link: '/ja/reference/architecture' },
            ]
          },
          {
            text: 'コミュニティ',
            items: [
              { text: '開発', link: '/ja/development' },
            ]
          }
        ],

        footer: {
          message: 'MITライセンスの下で公開されています。',
          copyright: 'Copyright © nbifrye'
        },

        editLink: {
          pattern: 'https://github.com/nbifrye/rmn/edit/main/docs/:path',
          text: 'GitHubでこのページを編集する'
        },

        outlineTitle: '目次',
        returnToTopLabel: 'ページの先頭へ',
        sidebarMenuLabel: 'メニュー',
        darkModeSwitchLabel: 'ダークモード',
        docFooter: {
          prev: '前のページ',
          next: '次のページ'
        }
      }
    }
  },

  themeConfig: {
    socialLinks: [
      { icon: 'github', link: 'https://github.com/nbifrye/rmn' }
    ],

    search: {
      provider: 'local'
    }
  }
})
