export type Dataset = {
  id: string;
  ownerType: "USER" | "ORGANIZATION";
  ownerId: string;
  title: string;
  slug: string;
  description: string;
  status: "ACTIVE" | "SUSPENDED" | "ARCHIVED";
  defaultVersionId: string;
  revision: number;
  createdAt: string;
  updatedAt: string;
};

type Envelope<T> = {
  code: string;
  message: string;
  data?: T;
};

export type DatasetListResult =
  | { ok: true; datasets: Dataset[] }
  | { ok: false; datasets: []; message: string };

const apiBaseUrl = process.env.API_BASE_URL ?? "http://localhost:8080";

export async function listDatasets(): Promise<DatasetListResult> {
  try {
    const response = await fetch(`${apiBaseUrl}/api/v1/datasets`, {
      cache: "no-store",
      headers: {
        "X-Owner-Type": process.env.DEV_OWNER_TYPE ?? "USER",
        "X-Owner-ID": process.env.DEV_OWNER_ID ?? "usr_local",
        "X-Actor-ID": process.env.DEV_ACTOR_ID ?? "usr_local",
      },
    });
    const payload = (await response.json()) as Envelope<Dataset[]>;
    if (!response.ok) {
      return { ok: false, datasets: [], message: payload.message || `API returned ${response.status}` };
    }
    return { ok: true, datasets: payload.data ?? [] };
  } catch {
    return {
      ok: false,
      datasets: [],
      message: "Dataset API is unavailable. Start the API service and refresh this page.",
    };
  }
}
