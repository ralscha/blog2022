Source code for my blog posts: https://blog.rasc.ch/

## Dependency compatibility and GPU upgrades

`pulumi-hetzner` stays on TypeScript 6 while the Pulumi SDK declares a `<7` peer
requirement. The CatBoost demo pins pandas 2.2, so its pandas-stubs updates stay
on the matching 2.2 line. Both restrictions are recorded in Dependabot settings.

The CUDA 13 PRs for `awsbatchgpu/catboost` and `awsbatchgpu/llm` need validation on
the actual GPU workers before deployment. Check the AWS Batch instance's driver
with `nvidia-smi`; CUDA 13.x needs driver branch R580 or newer, and CUDA 13.4's new
features require R615 or newer. Replacing a container image does not update the
host driver. See [NVIDIA's compatibility table](https://docs.nvidia.com/cuda/cuda-toolkit-release-notes/index.html).

For a migration, build the candidate images, run a small CatBoost training job
with GPU mode, and run the Pixtral sample on workers matching the production
instance type and GPU count. Check `torch.version.cuda` and the vLLM wheel's
CUDA/PyTorch requirements; changing the base image alone does not rebuild those
native extensions. The LLM demo pins an older vLLM release with CUDA 12-era
binaries, so validate a matching newer wheel or source build before using CUDA
13. Keep the previous image digest and worker configuration for rollback.
These GPU execution checks have not been run on this Windows workstation.
