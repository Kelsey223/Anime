import AnimeNotifier from "../AnimeNotifier"
import { requestIdleCallback } from "../Utils"

// Load
export function load(arn: AnimeNotifier, element: HTMLElement) {
	let url = element.dataset.url || (element as HTMLAnchorElement).getAttribute("href")

	if(!url) {
		arn.statusMessage.showError("Link doesn't have a target")
		return
	}

	arn.app.load(url)
}

// Diff
export async function diff(arn: AnimeNotifier, element: HTMLElement) {
	let url = element.dataset.url || (element as HTMLAnchorElement).getAttribute("href")

	if(!url) {
		arn.statusMessage.showError("Link doesn't have a target")
		return
	}

	try {
		await arn.diff(url)

		// Avoid instant layout thrashing
		requestIdleCallback(() => arn.scrollTo(element))
	} catch(err) {
		arn.statusMessage.showError(err)
	}
}