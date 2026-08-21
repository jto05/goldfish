// ModList displays the list of installed mods fetched from GET /api/mods.
// Implemented in Phase 5 — currently a placeholder.
function ModList() {
  return (
    <div className="rounded-lg border border-border h-full p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="font-semibold text-sm">Mods</h2>
        <button className="text-xs text-muted-foreground hover:text-foreground">+ Add</button>
      </div>
      <p className="text-sm text-muted-foreground">Coming soon.</p>
    </div>
  )
}

export default ModList
