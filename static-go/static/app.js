document.getElementById('check')?.addEventListener('click', async () => {
const out = document.getElementById('out')
out.textContent = 'Checking...'
try {
const r = await fetch('/healthz')
const txt = await r.text()
out.textContent = `status: ${r.status}\nbody: ${txt}`
} catch (err) {
out.textContent = `error: ${err}`
}
})