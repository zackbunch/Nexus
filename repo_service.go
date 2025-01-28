func (s *RepositoryService) GetDownloadURLs(ctx context.Context, repository, path string) ([]string, error) {
	if repository == "" {
		return nil, fmt.Errorf("repository name cannot be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	var (
		downloadURLs      []string
		continuationToken string
		totalAssets       int
	)

	// Configure rate limiting
	rateLimiter := time.NewTicker(time.Millisecond * 200) // 5 requests per second
	defer rateLimiter.Stop()

	spinner, _ := pterm.DefaultSpinner.WithText(fmt.Sprintf("Fetching download URLs from path '%s' in repository '%s'...", path, repository)).Start()

	for {
		select {
		case <-ctx.Done():
			spinner.Fail("Operation cancelled")
			return nil, ctx.Err()
		case <-rateLimiter.C:
			req := s.Client.restyClient.R().
				SetContext(ctx).
				SetQueryParams(map[string]string{
					"repository": repository,
					"q":          path,
				}).
				SetResult(&RepositoryServiceResponse{})

			if continuationToken != "" {
				req.SetQueryParam("continuationToken", continuationToken)
			}

			resp, err := req.Get("/service/rest/v1/search/assets")
			if err != nil {
				spinner.Fail("Failed to fetch download URLs from Nexus")
				return nil, fmt.Errorf("failed to fetch assets: %w", err)
			}

			if resp.IsError() {
				statusCode := resp.StatusCode()
				spinner.Fail(fmt.Sprintf("API error: %d", statusCode))
				return nil, fmt.Errorf("API error: %s (status code: %d)", resp.Status(), statusCode)
			}

			responseData, ok := resp.Result().(*RepositoryServiceResponse)
			if !ok {
				spinner.Fail("Failed to parse response")
				return nil, fmt.Errorf("failed to parse response data")
			}

			for _, asset := range responseData.Items {
				downloadURLs = append(downloadURLs, asset.DownloadURL)
			}
			totalAssets += len(responseData.Items)

			if responseData.ContinuationToken == "" {
				break
			}

			continuationToken = responseData.ContinuationToken
			spinner.UpdateText(fmt.Sprintf("Fetched %d assets so far...", totalAssets))
		}
	}

	spinner.Success(fmt.Sprintf("Fetched %d download URLs successfully.", totalAssets))

	if len(downloadURLs) == 0 {
		return nil, fmt.Errorf("no assets found in repository '%s' under path '%s'", repository, path)
	}

	return downloadURLs, nil
}