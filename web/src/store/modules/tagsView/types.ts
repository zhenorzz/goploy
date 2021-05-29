import { RouteLocationNormalizedLoaded } from 'vue-router'
export interface TagsViewState {
  visitedViews: RouteLocationNormalizedLoaded[]
  cachedViews: RouteLocationNormalizedLoaded[]
}
