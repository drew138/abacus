package rest

func mapImportResponse(added, updated, kept int) ImportResponse {
	return ImportResponse{Added: added, Updated: updated, Kept: kept}
}
