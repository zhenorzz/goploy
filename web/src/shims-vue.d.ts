declare module 'path-browserify' {
  import path from 'path'
  export default path
}

declare module 'virtual:svg-icons-names' {
  const svgIds: string[]
  export default svgIds
}
