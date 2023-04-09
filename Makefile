# Increases patch version. Pushes result tag.
.PHONY: release-patch
release-patch:
	@git tag $$(svu patch)
	@git push
	git push origin $$(svu next)

# Increases minor version. Pushes result tag.
.PHONY: release-minor
release-minor:
	@git tag $$(svu minor)
	@git push
	git push origin $$(svu next)

# Increases minor version. Pushes result tag.
.PHONY: release-major
release-major:
	@git tag $$(svu minor)
	@git push
	git push origin $$(svu next)

# Releases version based on commits. Pushes result tag.
.PHONY: release
release:
	@git tag $$(svu next)
	@git push
	git push origin $$(svu next)