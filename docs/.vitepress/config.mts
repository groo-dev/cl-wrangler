import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "cl",
  description: "Switch between multiple Cloudflare/Wrangler accounts",
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Installation', link: '/guide/installation' },
    ],

    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Installation', link: '/guide/installation' },
          { text: 'Usage', link: '/guide/usage' },
          { text: 'Configuration', link: '/guide/configuration' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/groo-dev/cl-wrangler' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 Groo Dev'
    },

    search: {
      provider: 'local'
    }
  }
})
