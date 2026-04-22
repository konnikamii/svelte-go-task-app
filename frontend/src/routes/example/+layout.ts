export const ssr = true

export const load = async ({ data }) => {
  return { layoutTs: 'layout.ts', ...data }
}
