import { describe, expect, it } from "vitest";
import { MISSING_REPO_DEPLOY_MESSAGE } from "./create-deploy-message";

describe("MISSING_REPO_DEPLOY_MESSAGE", () => {
  it("tells users how to recover from a git deploy without a repository", () => {
    expect(MISSING_REPO_DEPLOY_MESSAGE).toContain("Connect a repository in app settings");
    expect(MISSING_REPO_DEPLOY_MESSAGE).toContain("deploy a prebuilt image instead");
  });
});
