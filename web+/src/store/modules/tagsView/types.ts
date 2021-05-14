import { RouteLocationNormalizedLoaded } from 'vue-router'
export type TagsViewState = {
  visitedViews: RouteLocationNormalizedLoaded[]
  cachedViews: RouteLocationNormalizedLoaded[]
}
