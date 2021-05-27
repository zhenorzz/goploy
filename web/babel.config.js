module.exports = {
  presets: [
    '@vue/app',
    [
      '@babel/preset-env',
      {
        useBuiltIns: 'entry',
      },
    ],
  ],
  plugins: ['transform-class-properties'],
}
