package featureflag

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/miqdadyyy/fiber-featureflag/providers"
)

type FiberFeatureFlag struct {
	provider providers.IFeatureFlagProvider
}

func (f *FiberFeatureFlag) EnableFeatureFlag(ctx context.Context, key string) error {
	return f.provider.SetFeatureFlagStatus(ctx, key, true)
}

func (f *FiberFeatureFlag) DisableFeatureFlag(ctx context.Context, key string) error {
	return f.provider.SetFeatureFlagStatus(ctx, key, false)
}

func (f *FiberFeatureFlag) GetFeatureFlagStatus(ctx context.Context, key string) bool {
	return f.provider.GetFeatureFlagStatus(ctx, key)
}

func (f *FiberFeatureFlag) ToggleFeatureFlag(ctx context.Context, key string) error {
	return f.provider.SetFeatureFlagStatus(ctx, key, !f.provider.GetFeatureFlagStatus(ctx, key))
}

func (f *FiberFeatureFlag) GetAllFeatureFlags(ctx context.Context) (map[string]bool, error) {
	return f.provider.GetListOfFeatureFlags(ctx)
}

func (f *FiberFeatureFlag) PopulateFeatureFlag(ctx context.Context, keys []string) error {
	existsKeys, err := f.provider.GetListOfFeatureFlags(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, ok := existsKeys[key]; !ok {
			if err = f.provider.SetFeatureFlagStatus(ctx, key, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *FiberFeatureFlag) GetHandler(c *fiber.Ctx) error {
	if c.Get(fiber.HeaderAccept) == fiber.MIMEApplicationJSON {
		if c.Method() == fiber.MethodGet {
			// Check if a specific key is requested
			key := c.Query("key")
			if key != "" {
				status := f.provider.GetFeatureFlagStatus(c.Context(), key)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"key":     key,
					"enabled": status,
				})
			}

			// Return all feature flags
			featureFlags, err := f.provider.GetListOfFeatureFlags(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(fiber.StatusOK).JSON(featureFlags)
		}

		if c.Method() == fiber.MethodPost {
			var request struct {
				Key string `json:"key"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
			}

			if request.Key == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key is required"})
			}

			err := f.provider.SetFeatureFlagStatus(c.Context(), request.Key, true)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Feature flag enabled", "key": request.Key})
		}

		if c.Method() == fiber.MethodDelete {
			var request struct {
				Key string `json:"key"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
			}

			if request.Key == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key is required"})
			}

			err := f.provider.SetFeatureFlagStatus(c.Context(), request.Key, false)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Feature flag disabled", "key": request.Key})
		}

		if c.Method() == fiber.MethodPatch {
			var request struct {
				Key string `json:"key"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
			}

			if request.Key == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key is required"})
			}

			err := f.ToggleFeatureFlag(c.Context(), request.Key)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			// Get the updated status
			status := f.provider.GetFeatureFlagStatus(c.Context(), request.Key)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Feature flag toggled",
				"key":     request.Key,
				"enabled": status,
			})
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Status(fiber.StatusOK).SendString(IndexView)
}

func NewFiberFeatureFlag(provider providers.IFeatureFlagProvider) *FiberFeatureFlag {
	return &FiberFeatureFlag{
		provider: provider,
	}
}
