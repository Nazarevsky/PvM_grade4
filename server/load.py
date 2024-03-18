from locust import HttpUser, task


class WebUsr(HttpUser):
    @task
    def replace_round_brackets(self):
        self.client.post("/replace_round_brackets", json={
            "message": "(Hello world)"
        })
