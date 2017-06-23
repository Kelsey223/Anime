export function* findAll(className: string) {
	// getElementsByClassName failed for some reason.
	// TODO: Test getElementsByClassName again.
	let elements = document.querySelectorAll("." + className)
	
	for(let i = 0; i < elements.length; ++i) {
		yield elements[i] as HTMLElement
	}
}