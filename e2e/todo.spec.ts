import { test, expect } from "@playwright/test";

// ベースライン: タスクを追加 → 一覧に出る → 削除で消える。
// （完了トグルのリロード後永続化は既知の不具合のため検証しない）
test("add a task, see it listed, then delete it", async ({ page }) => {
  const title = `e2e task ${Date.now()}`;

  await page.goto("/");

  await page.getByTestId("new-task-input").fill(title);
  await page.getByTestId("add-button").click();

  const item = page.getByTestId("task-item").filter({ hasText: title });
  await expect(item).toHaveCount(1);

  await item.getByRole("button", { name: `delete ${title}` }).click();

  await expect(page.getByTestId("task-item").filter({ hasText: title })).toHaveCount(0);
});
